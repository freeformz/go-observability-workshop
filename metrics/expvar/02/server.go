package main

import (
	"errors"
	"expvar"
	"math/rand"
	"net/http"
	"os"
	"time"

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

func httpLoggingAndMetricsHandler(log logrus.FieldLogger, reqs, errs *expvar.Int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK // net/http returns 200 by default
		log = log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.String(),
		})
		defer func(t time.Time) {
			reqs.Add(1)
			log.WithField("status", status).WithField("duration", time.Since(t).Seconds()).Info()
		}(time.Now())

		if err := work(log); err != nil {
			status = http.StatusBadRequest
			errs.Add(1)
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

	// Expose the port value
	ep := expvar.NewString("Port")
	ep.Set(port)

	// Export the numbers
	reqs := expvar.NewInt("Requests")
	errs := expvar.NewInt("Errors")

	http.HandleFunc("/", httpLoggingAndMetricsHandler(log, reqs, errs))

	log.Info("Listening at: http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Errored with: " + err.Error())
	}
}
