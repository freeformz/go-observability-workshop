package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
)

const (
	serviceBURL = "http://localhost:8081"
)

func errorResponse(span *trace.Span, err error, w http.ResponseWriter) {
	span.SetStatus(trace.Status{Code: trace.StatusCodeInternal, Message: err.Error()})
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func queryServiceBHandler(c *http.Client, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := url // make a copy
		ctx, span := trace.StartSpan(r.Context(), "queryServiceBHandler")
		defer span.End()

		s := rand.Intn(99) + 1 // 1..100
		span.Annotate([]trace.Attribute{
			trace.Int64Attribute("s", int64(s)),
		}, "")

		// Pretend local computation before calling service b
		time.Sleep(time.Duration(s) * time.Millisecond / 4)
		span.SetStatus(trace.Status{Message: "local work complete"})

		if s <= 25 { // ~25% of the time call b's slow URL
			url = url + "/slow"
		}
		span.Annotate([]trace.Attribute{
			trace.StringAttribute("url", url),
		}, "")

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			errorResponse(span, errors.Wrap(err, "creating request"), w)
			return
		}

		resp, err := c.Do(req.WithContext(ctx))
		if err != nil {
			errorResponse(span, errors.Wrap(err, "doing request"), w)
			return
		}

		w.WriteHeader(resp.StatusCode)

		if resp.Body != nil {
			b, err := io.Copy(w, resp.Body)
			span.Annotate([]trace.Attribute{
				trace.Int64Attribute("proxied_bytes", b),
			}, "proxied")
			if err != nil {
				if b == 0 {
					errorResponse(span, errors.Wrap(err, "proxying bytes"), w)
				}
				return
			}
		}

		if resp.StatusCode/100 == 2 {
			w.Write([]byte(`a = :-)`))
		}
	}
}

func slowLocalWork(w http.ResponseWriter, r *http.Request) { // slow pretend work
	_, span := trace.StartSpan(r.Context(), "slowLocalWork")
	defer span.End()

	s := 100 + rand.Intn(200) // 100..300
	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("s", int64(s)),
	}, "")

	time.Sleep(time.Duration(s) * time.Millisecond)

	w.Write([]byte(`a = 🐢 `))
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})

	// curried log
	log := logrus.WithField("app", "servicea")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	je, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: "http://localhost:14268/api/traces",
		Process: jaeger.Process{
			ServiceName: "servicea",
			Tags: []jaeger.Tag{
				jaeger.StringTag("server", "1"), // could be hostname
				jaeger.StringTag("port", port),
				jaeger.BoolTag("demo", true),
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create the Jaeger exporter: %v", err)
	}
	trace.RegisterExporter(je)                                            //register the exporter
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}) // demo, so always sample

	// Expose the port value
	info := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "program_info",
		Help: "Info about the program.",
	},
		[]string{"port"},
	)
	prometheus.MustRegister(info)
	info.WithLabelValues(port).Set(1)

	durs := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request duration.",
			// Chosen because the range is 00-300 ms
			Buckets: []float64{.025, .05, .075, .1, .125, .15, .175, .2, .225, .250, .275, .300},
		},
		[]string{"handler", "code"},
	)
	prometheus.MustRegister(durs)
	durs.WithLabelValues("regularWork", strconv.Itoa(http.StatusOK))
	durs.WithLabelValues("regularWork", strconv.Itoa(http.StatusBadRequest))
	durs.WithLabelValues("slowWork", strconv.Itoa(http.StatusOK))
	durs.WithLabelValues("slowWork", strconv.Itoa(http.StatusBadRequest))

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	var oct ochttp.Transport
	c := http.Client{Transport: &oct, Timeout: 2 * time.Second} // always set sensible values for your service, never trust the defaults
	mux.Handle("/",
		ochttp.WithRouteTag(
			http.HandlerFunc(
				promhttp.InstrumentHandlerDuration(
					durs.MustCurryWith(prometheus.Labels{"handler": "queryServiceB"}),
					http.HandlerFunc(queryServiceBHandler(&c, serviceBURL)),
				),
			),
			"/",
		),
	)

	mux.Handle("/slow",
		ochttp.WithRouteTag(
			http.HandlerFunc(
				promhttp.InstrumentHandlerDuration(
					durs.MustCurryWith(prometheus.Labels{"handler": "slowLocalWork"}),
					http.HandlerFunc(slowLocalWork),
				),
			),
			"/slow",
		),
	)

	log.Info("Listening at: http://localhost:" + port)
	var pf b3.HTTPFormat
	if err := http.ListenAndServe(
		":"+port,
		&ochttp.Handler{
			Handler:     mux,
			Propagation: &pf,
		},
	); err != nil {
		log.Fatal("Errored with: " + err.Error())
	}
}
