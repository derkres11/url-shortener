package save

import (
	"errors"
	"modules/internal/config/lib/api/response"
	"modules/internal/config/lib/logger/sl"
	"modules/internal/config/lib/random"
	"modules/internal/config/storage"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	Response response.Response
	Alias    string `json:"alias,omitempty"`
}

const aliasLength = 8 // You can adjust the length as needed

//go:generate go run github.com/golang/mock/mockgen -destination=./mocks/url_saver_mock.go -package=mocks . URLSaver
type URLSaver interface {
	SaveURL(shortURL, originalURL string) error
}

func New(log *slog.Logger, saver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Save URL handler called")

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			// Failed to decode request body

			render.JSON(w, r, response.Error("Invalid request body"))

			return
		}

		log.Info("Request body decoded", "url", req.URL, "alias", req.Alias)

		if err := validator.New().Struct(req); err != nil {
			log.Error("Validation failed", sl.Err(err))

			render.JSON(w, r, response.Error("Invalid request body"))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		err = saver.SaveURL(alias, req.URL)
		if errors.Is(err, storage.ErrURLAlreadyExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("URL already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.String("alias", alias))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
