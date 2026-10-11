package middleware

import (
	"log/slog"
	"net/http"
	apperror "rps/internal/errors"
	"runtime/debug"

	"rps/internal/http/httpx"
)

func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error(
						"panic recovered",
						"panic", rec,
						"stack", string(debug.Stack()),
					)
					httpx.WriteErrorDetail(
						w,
						http.StatusInternalServerError,
						apperror.InternalServerErr.Error(),
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
