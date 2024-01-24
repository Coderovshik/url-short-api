package ping

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ping.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("ping request")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	}
}
