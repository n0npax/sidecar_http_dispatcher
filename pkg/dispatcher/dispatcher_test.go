package dispatcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mapping functions to vars to provide testing possibility
func sendRequest(ctx context.Context, r http.Handler, method string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, "/", nil)
	req.Header.Add("environment", "dev")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	return w
}

func TestDispatch(t *testing.T) {
	r, _, ctx := Router()

	w := sendRequest(ctx, r, "GET")
	assert.Equal(t, http.StatusOK, w.Code)
}

func init() { // nolint
	os.Setenv("SIDECAR_CONFIG", "../../config.yaml")
}
