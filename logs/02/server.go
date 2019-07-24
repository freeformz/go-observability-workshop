package main

import (
	"errors"
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

func httpLogginghandler(log logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK // net/http returns 200 by default
		log = log.WithFields(logrus.Fields{
			"method": r.Method,
			"path":   r.URL.String(),
		})
		defer func(t time.Time) {
			log.WithField("status", status).WithField("duration", time.Since(t).Seconds()).Info()
		}(time.Now())

		if err := work(log); err != nil {
			status = http.StatusBadRequest
			http.Error(w, "Nope", status)
			log.Error("OMG Error!")
			return
		}

		w.Write([]byte(`:-)`))
	}
}

func main() {
	/*logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
	})*/

	logrus.SetFormatter(&logrus.JSONFormatter{
		DisableColors: true,
	})

	// curried log
	log := logrus.WithField("app", "logs-02-server")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", httpLogginghandler(log))

	log.Info("Listening at: http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Errored with: " + err.Error())
	}
}
