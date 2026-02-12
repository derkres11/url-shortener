package save

import (
	"fmt"
	"modules/internal/config/lib/api/response"
	"modules/internal/config/lib/logger/sl"
	"net/http"

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

type URLSaver interface {
	SaveURL(shortURL, originalURL string) error
}

func New(log *slog.Logger, saver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Save URL handler called")

		var req Request
		err := render.DecodeJSON(r.body, &req)
		if err != nil {
			fmt.Errorf("Failed to decode request body: %w", err)

			render.JSON(w, r, response.Error("Invalid request body"))

			return
		}

		log.Info("Request body decoded", "url", req.URL, "alias", req.Alias)

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Validation failed", sl.Err(err))

			render.JSON(w, r, response.Error("Invalid request body"))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
	}
}
