package dispatcher

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/n0npax/sidecar_http_dispatcher/pkg/config"
)

var conf = config.GetConfig()
var client = &http.Client{}
var dispatchKey = conf.Key

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

	r.URL = u
	r.Host = u.Host
	r.RequestURI = ""
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
	resp, body := passRequest(req)
	return resp, body
}
