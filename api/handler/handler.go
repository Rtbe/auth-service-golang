package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

//Handler is a handler with nested router.
type Handler struct {
	Router  *chi.Mux
	Context context.Context
}

//New creates new Handler with nested router.
func New(ctx context.Context, router *chi.Mux) *Handler {
	return &Handler{
		Router:  router,
		Context: ctx,
	}
}

//RespondWithJSON is a helper for handling json responses.
func respondWithJSON(message string, payload interface{}, statusCode int, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonMap := make(map[string]interface{})
	jsonMap[message] = payload

	json.NewEncoder(w).Encode(jsonMap)
}

//RespondWithError is a helper for handling json responses with errors
func respondWithError(payload interface{}, statusCode int, w http.ResponseWriter) {
	respondWithJSON("error", payload, statusCode, w)
}
