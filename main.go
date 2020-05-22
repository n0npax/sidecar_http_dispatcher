package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	client := &http.Client{}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			panic("cannot create a req")
		}
		resp, _ := client.Do(req)
		w.Write([]byte(fmt.Sprintf("%v", resp)))
	})

	http.ListenAndServe(":3333", r)
}
