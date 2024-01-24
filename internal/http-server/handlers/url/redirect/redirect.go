package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Coderovshik/url-short/internal/lib/api/response"
	"github.com/Coderovshik/url-short/internal/lib/logger/slut"
	"github.com/Coderovshik/url-short/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func New(ctx context.Context, log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requeseId", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if len(alias) == 0 {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("empty alias"))

			return
		}

		url, err := urlGetter.GetURL(ctx, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("url not found", slut.ErrAttr(err))

				render.JSON(w, r, response.Error("url not found"))
			} else {
				log.Error("failed to get url")

				render.JSON(w, r, response.Error("failed to get url"))
			}

			return
		}

		log.Info("got url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
	}
}
