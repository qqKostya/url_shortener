package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url_shortener/internal/lib/api/response"
	"url_shortener/internal/lib/logger/sl"
	"url_shortener/internal/lib/random"
	"url_shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias, omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias, omitempty"`
}

// TODO: mb move to config
const aliasLength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.46.3 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.url.save.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("failed to validate request body", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			//TODO: checked in database
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("alias", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to save url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save url"))

			return
		}

		log.Info("url saved", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
