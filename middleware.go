package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const ctxKeyRequestID contextKey = "requestUUID"

func reqIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = assignRequestID(ctx)
		r = r.WithContext(ctx)

		slog.Info("Incoming request", "requestId", getRequestID(ctx), "method", r.Method, "requestURI", r.RequestURI, "remoteAddr", r.RemoteAddr)

		next.ServeHTTP(w, r)
	})
}

func assignRequestID(ctx context.Context) context.Context {
	reqID := uuid.New()
	return context.WithValue(ctx, ctxKeyRequestID, reqID.String())
}

func getRequestID(ctx context.Context) string {
	reqID := ctx.Value(ctxKeyRequestID)
	if id, ok := reqID.(string); ok {
		return id
	}
	return "undefined"
}
