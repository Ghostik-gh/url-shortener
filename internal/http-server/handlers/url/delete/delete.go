package delete

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

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
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

		log.Info("alias recieved", slog.Any("alias", alias))

		if err := urlDeleter.DeleteURL(alias); err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("alias not found", slog.String("alias", alias))
				render.JSON(w, r, resp.Error("alias not found"))
				return
			}
			log.Info("failed to delete alias")
			render.JSON(w, r, resp.Error("failed to delete alias"))
			return
		}
		log.Info("alias deleted", slog.String("alias", alias))
		render.JSON(w, r, resp.OK())
	}
}
