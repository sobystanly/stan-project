package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	rh *riskHandler
}

func NewHandler(rh *riskHandler) *Handler {
	return &Handler{rh: rh}
}

func NewRouter(h *Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(h.RequestIDMiddleware)
	for _, route := range h.GetRoutes() {
		hf := route.HandlerFunc
		router.Methods(route.Method).Name(route.Name).Handler(hf).Path(route.Pattern)
	}
	return router
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getQueryParam(key string, r *http.Request) string {
	return r.URL.Query().Get(key)
}
