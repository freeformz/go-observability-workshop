package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func work(log logrus.FieldLogger) error { // pretend work
	defer func(t time.Time) {
		log.WithField("work_seconds", time.Since(t).Seconds()).Info("Work complete")
	}(time.Now())

	s := rand.Intn(99) + 1 // 1..100
	time.Sleep(time.Duration(s) * time.Millisecond)

	var err error
	if s <= 25 { // ~25% of the time the work errors
		err = errors.New("OMG Error!")
	}
	return err
}

func slowWork(log logrus.FieldLogger) error { // slow pretend work
	s := 100 + rand.Intn(200) // 100..300
	defer log.WithField("work_seconds", float64(s)/1000).Info("Work complete")

	time.Sleep(time.Duration(s) * time.Millisecond)

	return nil
}

type workFunc func(logrus.FieldLogger) error

func httpLoggingAndMetricsHandler(log logrus.FieldLogger, durs prometheus.ObserverVec, wf workFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK // net/http returns 200 by default
		log = log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.String(),
		})
		defer func(t time.Time) {
			secs := time.Since(t).Seconds()
			durs.WithLabelValues(strconv.Itoa(status)).Observe(secs)
			log.WithField("status", status).WithField("duration", secs).Info()
		}(time.Now())

		if err := wf(log); err != nil {
			status = http.StatusBadRequest
			http.Error(w, "Nope", status)
			log.Error("OMG Error!")
			return
		}

		w.Write([]byte(`:-)`))
	}
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})

	// curried log
	log := logrus.WithField("app", "logs-02-server")

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

	http.HandleFunc("/", httpLoggingAndMetricsHandler(
		log,
		durs.MustCurryWith(prometheus.Labels{"handler": "regularWork"}),
		work,
	))

	http.HandleFunc("/slow", httpLoggingAndMetricsHandler(
		log,
		durs.MustCurryWith(prometheus.Labels{"handler": "slowWork"}),
		slowWork,
	))

	log.Info("Listening at: http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Errored with: " + err.Error())
	}
}
