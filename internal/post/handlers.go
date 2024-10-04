package post

import (
	"blog_api/internal/storage"
	"blog_api/pkg/logger"
	"blog_api/pkg/utils"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type CreatePostRequest struct{
	Title string `json:"title" validate:"required"`
	Body string `json:"body" validate:"required"`
	Author string `json:"author" validate:"required"`
}

type CreatePostResponse struct {
	utils.Response

}

func CreatePost(log *slog.Logger, queryTool storage.PostQueryFunctions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req CreatePostRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, utils.Error("error decoding request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validErr := err.(validator.ValidationErrors)

			log.Error("Validate error", logger.Err(validErr))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, validErr)
			return
		}
		post := storage.Post{
			Title: req.Title,
			Body: req.Body,
			Author: req.Author,
		}

		id, err := queryTool.NewPost(context.Background(), post)
		if err != nil {
			log.Error("Error while adding new post", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, utils.Error("db error"))
			return
		}
		log.Info("Inserted post", slog.Int64("id", int64(id)))

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, utils.OK())
	}
}