package save

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Coderovshik/url-short/internal/lib/api/response"
	"github.com/Coderovshik/url-short/internal/lib/logger/slut"
	"github.com/Coderovshik/url-short/internal/lib/random"
	"github.com/Coderovshik/url-short/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config
const aliasLength = 6

type URLSaver interface {
	SaveURL(ctx context.Context, urlToSave string, alias string) error
	ExistAlias(ctx context.Context, alias string) bool
}

func New(ctx context.Context, log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
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

		alias := req.Alias
		if len(alias) == 0 {
			alias = random.NewRandomString(aliasLength)
			for urlSaver.ExistAlias(ctx, alias) {
				alias = random.NewRandomString(aliasLength)
			}
		}

		err = urlSaver.SaveURL(ctx, req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))

				render.JSON(w, r, response.Error("url alredy exists"))
			} else {
				log.Error("failed to add url", slut.ErrAttr(err))

				render.JSON(w, r, response.Error("failed to add url"))
			}

			return
		}

		log.Info("url added", slog.String("alias", alias))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
