package redirect

import (
	"errors"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

// URLGetter is an interface for getting url by alias.
//
//go:generate mockery --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("alias is empty"))
			return
		}

		log.Info("alias recieved", slog.Any("request", alias))

		resUrl, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", slog.String("alias", alias))
				render.JSON(w, r, resp.Error("url not found"))
				return
			}
			log.Info("failed to get url")
			render.JSON(w, r, resp.Error("failed to get url"))
			return
		}

		log.Info("url recieved", slog.String("resUrl", resUrl))

		http.Redirect(w, r, resUrl, http.StatusFound)

	}
}
