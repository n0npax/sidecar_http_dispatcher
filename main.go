package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
		u, err := url.Parse("http://example.com")
		if err != nil {
			log.Fatal(err)
		}

		r.URL = u
		r.Host = u.Host
		r.RequestURI = ""

		resp, err := client.Do(r)
		if err != nil {
			panic(err)
		}
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
