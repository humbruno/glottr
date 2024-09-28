package main

import "net/http"

type handlers struct{}

func (h *handlers) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Yep"))
}
