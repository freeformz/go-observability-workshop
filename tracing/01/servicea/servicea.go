package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func errorResponse(log logrus.FieldLogger, err error, w http.ResponseWriter, status int) {
	log.Error(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func queryServiceBHandler(c *http.Client, log logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		log = log.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.String(),
			"request_id": id,
		})
		status := http.StatusOK // net/http returns 200 by default
		defer func(t time.Time) {
			log.WithField("status", status).WithField("duration", time.Since(t).Seconds()).Info()
		}(time.Now())

		s := rand.Intn(99) + 1 // 1..100
		// Pretend local computation before calling service b
		time.Sleep(time.Duration(s) * time.Millisecond / 4)

		url := "http://localhost:8081"
		if s <= 25 { // ~25% of the time call the slow URL
			url = url + "/slow"
		}
		log = log.WithField("url", url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			status = http.StatusInternalServerError
			errorResponse(log, errors.Wrap(err, "creating request"), w, status)
			return
		}
		req.Header.Set("X-Request-ID", id)

		resp, err := c.Do(req)
		if err != nil {
			status = http.StatusInternalServerError
			errorResponse(log, errors.Wrap(err, "doing request"), w, status)
			return
		}
		w.WriteHeader(resp.StatusCode)

		if resp.Body != nil {
			b, err := io.Copy(w, resp.Body)
			log.WithField("proxied_bytes", b).Info()
			if err != nil {
				status = http.StatusInternalServerError
				if b == 0 {
					errorResponse(log, errors.Wrap(err, "proxying bytes"), w, status)
				}
				return
			}
		}

		if resp.StatusCode/100 == 2 {
			w.Write([]byte(`a = :-)`))
		}
	}
}

func slowHandler(log logrus.FieldLogger) http.HandlerFunc { // slow pretend work
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		log = log.WithFields(logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.String(),
			"request_id": id,
		})
		defer func(t time.Time) {
			log.WithField("status", http.StatusOK).WithField("duration", time.Since(t).Seconds()).Info()
		}(time.Now())

		s := 100 + rand.Intn(200) // 100..300
		time.Sleep(time.Duration(s) * time.Millisecond)

		w.Write([]byte("a = 🐢 "))
	}
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

	http.Handle("/metrics", promhttp.Handler())

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

	var c http.Client
	c.Timeout = 2 * time.Second // always set sensible values for your service, never trust the defaults
	http.HandleFunc("/", promhttp.InstrumentHandlerDuration(
		durs.MustCurryWith(prometheus.Labels{"handler": "regularWork"}),
		http.HandlerFunc(queryServiceBHandler(&c, log)),
	))

	http.HandleFunc("/slow", promhttp.InstrumentHandlerDuration(
		durs.MustCurryWith(prometheus.Labels{"handler": "slowWork"}),
		http.HandlerFunc(slowHandler(log)),
	))

	log.Info("Listening at: http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Errored with: " + err.Error())
	}
}
