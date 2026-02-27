package redirect

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

type URLGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("Alias is empty")
			http.Error(w, "Alias is required", http.StatusBadRequest)
			return
		}

		resUrl, err = urlGetter.GetUrl(alias)
		if err != nil {
			log.Error("Failed to get URL", "error", err)
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		log.Info("Redirecting to URL", "url", resUrl)

		http.Redirect(w, r, resUrl, http.StatusFound)

	}
}
