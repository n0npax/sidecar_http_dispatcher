package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/valve"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/n0npax/sidecar_http_dispatcher/pkg/dispatcher"
	"github.com/n0npax/sidecar_http_dispatcher/pkg/utils"
)

const shutdownTimeout = 5 * time.Second

func main() {
	valv := valve.New()
	baseCtx := valv.Context()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.HandleFunc("/*", handleAndPass)

	addr := fmt.Sprintf(":%s", utils.GetEnv("SIDECAR_PORT", "5000"))
	log.Printf("Staring server on address: %s", addr)

	srv := http.Server{Addr: addr, Handler: chi.ServerBaseContext(baseCtx, r)}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Println("shutting down..")

			// send valv
			if err := valv.Shutdown(shutdownTimeout * time.Second); err != nil {
				log.Fatal(err)
			}

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel() // nolint

			// graceful shutdown
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}

			select {
			case <-time.After(shutdownTimeout + time.Second):
				fmt.Println("not all connections done. Killing anyway via cancel()")
			case <-ctx.Done():
			}
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handleAndPass(w http.ResponseWriter, r *http.Request) {
	if err := valve.Lever(r.Context()).Open(); err != nil {
		panic(err)
	}
	defer valve.Lever(r.Context()).Close()

	resp, body := dispatcher.Dispatch(r)
	if _, err := w.Write(body); err != nil {
		panic(err)
	}

	for k, v := range resp.Header {
		w.Header().Add(k, v[0])
	}

	defer resp.Body.Close()
}
