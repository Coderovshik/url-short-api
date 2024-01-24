package delete

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Coderovshik/url-short/internal/lib/api/response"
	"github.com/Coderovshik/url-short/internal/lib/logger/slut"
	"github.com/Coderovshik/url-short/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

type URLDeleter interface {
	DeleteURL(ctx context.Context, alias string) error
}

func New(ctx context.Context, log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			msg := "failed to decode request body"
			log.Error(msg, slut.ErrAttr(err))
			render.JSON(w, r, response.Error(msg))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		v := validator.New()
		if err := v.Struct(req); err != nil {
			var errs validator.ValidationErrors
			if errors.As(err, &errs) {
				msg := "invalid request"
				log.Error(msg, slut.ErrAttr(err))
				render.JSON(w, r, response.ValidationError(errs))
			}

			return
		}

		err = urlDeleter.DeleteURL(ctx, req.Alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("url not found", slut.ErrAttr(err))

				render.JSON(w, r, response.Error("url not found"))
			} else {
				log.Error("failed to delete url", slut.ErrAttr(err))

				render.JSON(w, r, response.Error("failed to delte url"))
			}

			return
		}

		log.Info("alias delted", slog.String("alias", req.Alias))

		render.JSON(w, r, response.OK())
	}
}
