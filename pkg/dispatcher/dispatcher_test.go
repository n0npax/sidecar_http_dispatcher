package dispatcher

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mapping functions to vars to provide testing possibility
func sendRequest(ctx context.Context, r http.Handler, method, path, header string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Add("environment", header)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestDispatch(t *testing.T) {
	testCases := []struct {
		code   int
		path   string
		header string
	}{
		{path: "/", code: 200, header: "dev"},
		{path: "/", code: 500, header: "qa"},
		{path: "/foo", code: 404, header: "dev"},
	}
	for _, test := range testCases {
		path, header, code := test.path, test.header, test.code
		t.Run(fmt.Sprintf("%v:%v", path, header), func(t *testing.T) {
			r, _, ctx := Router()
			w := sendRequest(ctx, r, "GET", path, header)
			assert.Equal(t, w.Code, code)
		})
	}
}

func init() { // nolint
	os.Setenv("SIDECAR_CONFIG", "../../config.yaml")
}
