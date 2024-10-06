package main

import (
	"encoding/json"
	"net/http"
)

const one_mb = 1_048_578

type errEnvelope struct {
	Error string `json:"error"`
}

type dataEnvelope struct {
	Data any `json:"data"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	return writeJSON(w, status, &errEnvelope{Error: message})
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := one_mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	return writeJSON(w, status, &dataEnvelope{Data: data})
}
