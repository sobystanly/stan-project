package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"stan-project/data"
	"testing"
)

func TestNewRiskHandler(t *testing.T) {
	t.Run("successfully initialize risk handler", func(t *testing.T) {
		mrl := &mockRiskLogic{}
		actual := NewRiskHandler(mrl)
		assert.Equal(t, &riskHandler{riskLogic: mrl}, actual)
	})
}

func TestRiskHandler_Add(t *testing.T) {
	t.Run("successfully add a new risk", func(t *testing.T) {
		h := NewRiskHandler(&mockRiskLogic{
			risk: data.Risk{
				ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
				Title:       "threat 1",
				Description: "DDOS threat",
				State:       "open",
			},
		})

		expected := data.Risk{
			ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "open",
		}

		req, err := http.NewRequest(http.MethodPost, "/v1/risks", bytes.NewBuffer(getTestData()))
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(req.Context(), "requestID", requestID)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.Add(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp data.Risk
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("error decoding response: %s", err)
		}

		assert.Equal(t, expected, resp)
	})

	t.Run("failed to add a new risk, error from logic", func(t *testing.T) {
		h := NewRiskHandler(&mockRiskLogic{
			err: errors.New("some error"),
		})

		req, err := http.NewRequest(http.MethodPost, "/v1/risks", bytes.NewBuffer(getTestData()))
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(req.Context(), "requestID", requestID)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.Add(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("failed to add a new risk, invalid request", func(t *testing.T) {
		h := NewRiskHandler(&mockRiskLogic{})

		req, err := http.NewRequest(http.MethodPost, "/v1/risks", bytes.NewBuffer([]byte(`{`)))
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(req.Context(), "requestID", requestID)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.Add(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRiskHandler_GetByID(t *testing.T) {

	t.Run("successfully fetch risk by ID", func(t *testing.T) {
		h := NewRiskHandler(&mockRiskLogic{
			risk: data.Risk{
				ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
				Title:       "threat 1",
				Description: "DDOS threat",
				State:       "open",
			},
		})

		expected := data.Risk{
			ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
			Title:       "threat 1",
			Description: "DDOS threat",
			State:       "open",
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/risks/%s", "c7041e22-15c1-4293-9b43-c54c8dd4b909"), nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": "c7041e22-15c1-4293-9b43-c54c8dd4b909"})

		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(req.Context(), "requestID", requestID)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp data.Risk
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("error decoding response: %s", err)
		}

		assert.Equal(t, expected, resp)
	})

	t.Run("failed to fetch risk by ID, error from logic", func(t *testing.T) {
		h := NewRiskHandler(&mockRiskLogic{
			err: errors.New("some error from logic"),
		})

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/risks/%s", "c7041e22-15c1-4293-9b43-c54c8dd4b909"), nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": "c7041e22-15c1-4293-9b43-c54c8dd4b909"})

		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(req.Context(), "requestID", requestID)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("failed to fetch risk by ID, invalid ID", func(t *testing.T) {
		h := NewRiskHandler(&mockRiskLogic{
			err: errors.New("some error from logic"),
		})

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/risks/%s", "c7041e22-15c"), nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		req = mux.SetURLVars(req, map[string]string{"id": "c7041e22-15c"})

		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(req.Context(), "requestID", requestID)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRiskHandler_GetAll(t *testing.T) {
	t.Run("successfully get all risks", func(t *testing.T) {
		h := NewRiskHandler(&mockRiskLogic{
			paginatedRisk: data.PaginatedResponse{
				TotalCount: 1,
				Risks: []data.Risk{
					{
						ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
						Title:       "threat 1",
						Description: "DDOS threat",
						State:       "open",
					},
				},
			},
		})

		expected := data.PaginatedResponse{
			TotalCount: 1,
			Risks: []data.Risk{
				{
					ID:          uuid.MustParse("c7041e22-15c1-4293-9b43-c54c8dd4b909"),
					Title:       "threat 1",
					Description: "DDOS threat",
					State:       "open",
				},
			},
		}

		req, err := http.NewRequest(http.MethodGet, "/v1/risks", nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(req.Context(), "requestID", requestID)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.GetAll(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp data.PaginatedResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("error decoding response: %s", err)
		}

		assert.Equal(t, expected, resp)

	})

	t.Run("failed to get all risks, error from logic", func(t *testing.T) {
		h := NewRiskHandler(&mockRiskLogic{
			err: errors.New("some error"),
		})

		req, err := http.NewRequest(http.MethodGet, "/v1/risks", nil)
		if err != nil {
			t.Fatalf("error creating request: %s", err)
		}

		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(req.Context(), "requestID", requestID)

		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		h.GetAll(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})
}

func getTestData() []byte {
	return []byte(`
					{
						"title": "threat 1",
						"description": "DDOS threat",
						"state": "open"
					}
				`)
}

type mockRiskLogic struct {
	risk          data.Risk
	err           error
	paginatedRisk data.PaginatedResponse
}

func (m mockRiskLogic) Add(ctx context.Context, risk data.Risk) (data.Risk, error) {
	return m.risk, m.err
}

func (m mockRiskLogic) GetByID(ctx context.Context, ID uuid.UUID) (data.Risk, error) {
	return m.risk, m.err
}

func (m mockRiskLogic) GetAll(ctx context.Context, options data.Options) (data.PaginatedResponse, error) {
	return m.paginatedRisk, m.err
}
