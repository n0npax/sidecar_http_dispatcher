package dispatcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/valve"
	"github.com/n0npax/sidecar_http_dispatcher/pkg/config"
)

// mapping functions to vars to provide testing possibility
var conf config.Config      // nolint
var client = &http.Client{} // nolint
var dispatchKey string      // nolint

func patch(r *http.Request) *http.Request {
	dk := r.Header.Get(dispatchKey)
	destination := conf.Destination

	for k, v := range conf.Rewrites {
		if k == dk {
			destination = v.Destination
		}
	}

	u, err := url.Parse(destination)
	if err != nil {
		log.Fatal(err)
	}

	r.RequestURI, r.URL.Scheme, r.Host = "", u.Scheme, u.Host
	r.URL.Host = r.Host

	return r
}

func passRequest(r *http.Request) (*http.Response, []byte, int) {
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	return resp, body, resp.StatusCode
}

func dispatch(r *http.Request) (*http.Response, []byte, int) {
	req := patch(r)
	return passRequest(req)
}

func handleAndPass(w http.ResponseWriter, r *http.Request) {
	if err := valve.Lever(r.Context()).Open(); err != nil {
		panic(err)
	}
	defer valve.Lever(r.Context()).Close()

	resp, body, code := dispatch(r)
	if _, err := w.Write(body); err != nil {
		panic(err)
	}

	for k, v := range resp.Header {
		w.Header().Add(k, v[0])
	}

	fmt.Printf("RRRRRRRRRr %v <<<", code)
	w.WriteHeader(code)
	w.Header().Add("HTTP/1.1", resp.Status)

	defer resp.Body.Close()
}

func Router() (*chi.Mux, *valve.Valve, context.Context) {
	conf = config.GetConfig()
	dispatchKey = conf.Key
	valv := valve.New()
	baseCtx := valv.Context()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.HandleFunc("/*", handleAndPass)

	return r, valv, baseCtx
}
