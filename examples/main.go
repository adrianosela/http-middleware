package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adrianosela/middleman"
	"github.com/adrianosela/sslmgr"
	"github.com/gorilla/mux"
)

func main() {

	h := middleman.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Im alive!"))
	}))

	rtr := mux.NewRouter()
	rtr.Methods(http.MethodGet).Path("/healthcheck").Handler(h)

	ss, err := sslmgr.NewSecureServer(sslmgr.ServerConfig{
		Hostnames: []string{os.Getenv("CN_FOR_CERTIFICATE")},
		Handler:   rtr,
		ServeSSLFunc: func() bool {
			return false
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	ss.ListenAndServe()
}
