package main

import (
	"log"
	"net/http"

	"github.com/adrianosela/middleman"
)

func main() {
	logger := middleman.NewLogger(true, true, true)

	functional := middleman.NewFunctional(func() {
		log.Println("doing X before the request was received")
	}, func() {
		log.Println("doing X after the request was fulfilled")
	})

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello World!"))
	})

	mm := middleman.NewMiddleman(logger, functional)

	http.ListenAndServe(":80", mm.Wrap(h))
}
