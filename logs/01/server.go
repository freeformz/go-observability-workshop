package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

func work() error { // pretend work
	defer func(t time.Time) {
		log.Printf("Work took %2.3fs\n", time.Since(t).Seconds())
	}(time.Now())

	s := rand.Intn(99) + 1 // 1..100
	time.Sleep(time.Duration(s) * time.Millisecond)

	var err error
	if s <= 25 { // ~25% of the time the work errors
		err = errors.New("OMG Error!")
	}
	return err
}

func handler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK // net/http returns 200 by default
	defer func(t time.Time) {
		log.Printf("%s %q => %d (%2.3fs)\n", r.Method, r.URL.String(), status, time.Since(t).Seconds())
	}(time.Now())

	if err := work(); err != nil {
		status = http.StatusBadRequest
		http.Error(w, ":-(", status)
		log.Println("Error:", err.Error())
		return
	}

	w.Write([]byte(`:-)`))
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler)

	log.Println("Listening at: http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Errored with: " + err.Error())
	}
}
