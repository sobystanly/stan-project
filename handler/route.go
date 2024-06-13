package handler

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func (h *Handler) GetRoutes() []Route {
	return []Route{
		//Health check endpoint
		{
			Name:        "CheckHealth",
			Method:      http.MethodGet,
			Pattern:     "/risks/health",
			HandlerFunc: h.CheckHealth,
		},

		//Risk endpoints
		{
			Name:        "Create a Risk",
			Method:      http.MethodPost,
			Pattern:     "/v1/risks",
			HandlerFunc: h.rh.Add,
		},
		{
			Name:        "Get a Risk By ID",
			Method:      http.MethodGet,
			Pattern:     "/v1/risks/{id}",
			HandlerFunc: h.rh.GetByID,
		},
		{
			Name:        "Get All Risks",
			Method:      http.MethodGet,
			Pattern:     "/v1/risks",
			HandlerFunc: h.rh.GetAll,
		},
	}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
	return
}

// RequestIDMiddleware generate and add a requestID to each request
func (h *Handler) RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique requestID using Google's UUID library
		requestID := uuid.New().String()

		// Add the requestID to the request context
		ctx := context.WithValue(r.Context(), "requestID", requestID)

		// Add the requestID as a header in the response
		w.Header().Set("X-Request-ID", requestID)

		// Call the next handler with the modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
