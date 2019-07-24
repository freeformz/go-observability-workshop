package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
)

func workHandler(w http.ResponseWriter, r *http.Request) { // pretend work
	_, span := trace.StartSpan(r.Context(), "workHandler")
	defer span.End()

	s := rand.Intn(99) + 1 // 1..100
	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("s", int64(s)),
	}, "")

	time.Sleep(time.Duration(s) * time.Millisecond)

	switch {
	case s <= 25:
		http.Error(w, "cache miss", http.StatusBadRequest)
	default:
		w.Write([]byte(`b = :-) `))
	}
}

func slowWorkHandler(w http.ResponseWriter, r *http.Request) { // slow pretend work
	_, span := trace.StartSpan(r.Context(), "slowWorkHandler")
	defer span.End()

	s := 100 + rand.Intn(200) // 100..300
	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("s", int64(s)),
	}, "")
	time.Sleep(time.Duration(s) * time.Millisecond)

	w.Write([]byte(`b = 🐢 `))
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})

	// curried log
	log := logrus.WithField("app", "serviceb")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	je, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: "http://localhost:14268/api/traces",
		Process: jaeger.Process{
			ServiceName: "serviceb",
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

	mux.Handle("/",
		ochttp.WithRouteTag(
			http.HandlerFunc(
				promhttp.InstrumentHandlerDuration(
					durs.MustCurryWith(prometheus.Labels{"handler": "regularWork"}),
					http.HandlerFunc(workHandler),
				),
			),
			"/",
		),
	)

	mux.Handle("/slow",
		ochttp.WithRouteTag(
			http.HandlerFunc(
				promhttp.InstrumentHandlerDuration(
					durs.MustCurryWith(prometheus.Labels{"handler": "slowWork"}),
					http.HandlerFunc(slowWorkHandler),
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
