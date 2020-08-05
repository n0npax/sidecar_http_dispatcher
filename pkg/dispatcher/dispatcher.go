package dispatcher

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/n0npax/sidecar_http_dispatcher/pkg/config"
)

// mapping functions to vars to provide testing possibility.
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

	r.RequestURI, r.URL.Scheme, r.Host, r.URL.Host = "", u.Scheme, u.Host, u.Host

	return r
}

func passRequest(r *http.Request) (*http.Response, []byte, int) {
	resp, err := client.Do(r)
	if err != nil {
		panic(err) // better handling required
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
	resp, body, read := passRequest(req)

	defer resp.Body.Close()

	return resp, body, read
}

func handleAndPass(c *gin.Context) {
	r := c.Request
	updateXFF(r)
	resp, body, code := dispatch(r)

	defer resp.Body.Close()

	if err := c.ShouldBindHeader(resp.Header); err != nil {
		log.Println(err)
	}

	c.String(code, string(body))
}

func Router() *gin.Engine {
	conf = config.GetConfig()
	dispatchKey = conf.Key

	r := gin.Default()

	r.Any("/*path", handleAndPass)

	return r
}

func updateXFF(r *http.Request) {
	if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if prior, ok := r.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}

		r.Header.Set("X-Forwarded-For", clientIP)
	}
}
