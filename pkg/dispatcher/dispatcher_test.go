package dispatcher

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mapping functions to vars to provide testing possibility
func sendRequest(r http.Handler, method string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, "/", nil)
	req.Header.Add("environment", "dev")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	return w
}

func TestDispatch(t *testing.T) {
	_, _, r := Dispatcher()
	w := sendRequest(r, "GET")
	assert.Equal(t, http.StatusOK, w.Code)
}

func init() { // nolint
	os.Setenv("SIDECAR_CONFIG", "../../config.yaml")
}
