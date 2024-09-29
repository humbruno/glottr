package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type apiError struct {
	Msg        any `json:"msg"`
	StatusCode int `json:"statusCode"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func (e apiError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func newApiError(statusCode int, err error) apiError {
	return apiError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func invalidRequestData(errors map[string]string) apiError {
	return apiError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}

func invalidJSON() apiError {
	return newApiError(http.StatusBadRequest, fmt.Errorf("invalid JSON request data"))
}

func makeHandlerFunc(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiErr, ok := err.(apiError); ok {
				writeJSON(w, apiErr.StatusCode, apiErr)
			} else {
				errResp := map[string]any{
					"statusCode": http.StatusInternalServerError,
					"msg":        "internal server error",
				}
				writeJSON(w, http.StatusInternalServerError, errResp)
			}
			slog.Error("HTTP API error", "error", err.Error(), "path", r.URL.Path)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Error writing JSON response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
