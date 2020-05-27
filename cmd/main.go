package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/n0npax/sidecar_http_dispatcher/pkg/dispatcher"
	"github.com/n0npax/sidecar_http_dispatcher/pkg/utils"
)

const shutdownTimeout = 5 * time.Second

func main() {
	r := dispatcher.Router()

	addr := fmt.Sprintf(":%s", utils.GetEnv("SIDECAR_PORT", "5000"))
	log.Printf("Staring server on address: %s", addr)

	srv := http.Server{Addr: addr, Handler: r}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
