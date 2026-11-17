package url

import (
	resp "awesomeProject/Internal/lib/api/response"
	"awesomeProject/Internal/lib/random"
	"awesomeProject/Internal/lib/sl"
	"awesomeProject/storage"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias" validate:"required,url"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}
type UrlSaver interface {
	SaveUrl(saveUrl string, alias string, ctx context.Context) (int64, error)
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver UrlSaver, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/url/save.go/New"

		log = log.With(
			slog.String("op", op),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("Error decoded  request"))

			return
		}

		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("Invalid request:", sl.Err(err))
			validatorErr := err.(validator.ValidationErrors)

			render.JSON(w, r, resp.ValidationErrors(validatorErr))

			return

		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		id, err := urlSaver.SaveUrl(req.Url, alias, r.Context())
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("Url already exists", slog.String("url", req.Url))
			render.JSON(w, r, resp.Error("Url already exists"))

			return
		}
		if err != nil {
			log.Error("failed add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added ", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}

}
