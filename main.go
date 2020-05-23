package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// TODO setup nice params
	client := &http.Client{}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			panic(err)
		}
		req.Header.Add("sidecar_http_dispatched", "true")
		for k, v := range r.Header {
			fmt.Printf("-> %v %v\n", k, v[0])
			req.Header.Add(k, v[0])
		}

		resp, _ := client.Do(req)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		w.Write(body)

		for k, v := range resp.Header {
			fmt.Printf("<- %v %v\n", k, v[0])
			w.Header().Add(k, v[0])
		}
	})

	http.ListenAndServe(getenv("SIDECAR_PORT", ":5000"), r)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
