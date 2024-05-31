package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CheckHealth(t *testing.T) {
	t.Run("Successfully return as healthy", func(t *testing.T) {
		h := NewHandler(&riskHandler{})

		req, err := http.NewRequest(http.MethodGet, "/risks/health", nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		w := httptest.NewRecorder()

		h.CheckHealth(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestNewRouter(t *testing.T) {
	t.Run("successfully initialise http router", func(t *testing.T) {
		h := NewHandler(&riskHandler{})
		router := NewRouter(h)
		assert.NotNil(t, router)
	})
}
