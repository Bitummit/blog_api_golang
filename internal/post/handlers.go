package post

import (
	"blog_api/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type CreatePostRequest struct{
	title string `json:"title" validate:"required"`
	body string `json:"body" validate:"required"`
	author string `json:"author" validate:"required"`
}

type CreatePostResponse struct {
	Status string
	Error string

}


func CreatePost(log *slog.Logger, storage *storage.PostQueryFunctions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)


	}
}