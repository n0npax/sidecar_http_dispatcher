package dispatcher_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/n0npax/sidecar_http_dispatcher/pkg/dispatcher"
	"github.com/stretchr/testify/assert"
)

// mapping functions to vars to provide testing possibility.
func sendRequest(r http.Handler, method, path, header string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	req.Header.Add("environment", header)

	r.ServeHTTP(w, req)

	return w
}

func TestDispatch(t *testing.T) {
	testCases := []struct {
		code   int
		path   string
		header string
	}{
		{path: "/foo", code: http.StatusNotFound, header: "dev"},
		{path: "/", code: http.StatusOK, header: "dev"},
		{path: "/", code: http.StatusInternalServerError, header: "qa"},
	}
	for _, test := range testCases {
		path, header, code := test.path, test.header, test.code
		t.Run(fmt.Sprintf("%s:%s", path, header), func(t *testing.T) {
			r := dispatcher.Router()
			w := sendRequest(r, "GET", path, header)
			assert.Equal(t, code, w.Code)
		})
	}
}

func init() { // nolint
	os.Setenv("SIDECAR_CONFIG", "../../config.yaml")
}
