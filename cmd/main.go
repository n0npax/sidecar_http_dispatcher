package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/n0npax/sidecar_http_dispatcher/pkg/dispatcher"
	"github.com/n0npax/sidecar_http_dispatcher/pkg/utils"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	//r.Use(middleware.Profiler)
	r.Use(middleware.Recoverer)

	r.HandleFunc("/*", handleAndPass)

	addr := fmt.Sprintf(":%s", utils.GetEnv("SIDECAR_PORT", "5000"))
	log.Printf("Staring server on address: %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func handleAndPass(w http.ResponseWriter, r *http.Request) {
	resp, body := dispatcher.Dispatch(r)
	if _, err := w.Write(body); err != nil {
		panic(err)
	}

	for k, v := range resp.Header {
		w.Header().Add(k, v[0])
	}

	defer resp.Body.Close()
}
