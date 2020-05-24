package dispatcher

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/n0npax/sidecar_http_dispatcher/pkg/config"
)

// mapping functions to vars to provide testing possibility
var conf = config.GetConfig() // nolint
var client = &http.Client{}   // nolint
var dispatchKey = conf.Key    // nolint

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

	r.RequestURI, r.URL, r.Host = "", u, u.Host

	return r
}

func passRequest(r *http.Request) (*http.Response, []byte) {
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	return resp, body
}

func Dispatch(r *http.Request) (*http.Response, []byte) {
	req := patch(r)
	return passRequest(req)
}
