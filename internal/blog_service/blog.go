package blogservice

import (
	"github.com/Bitummit/blog_api_golang/internal"
	"github.com/Bitummit/blog_api_golang/pkg/logger"
	"github.com/Bitummit/blog_api_golang/pkg/utils"
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)


type PostQueryFunctions interface {
	NewPost(context.Context, internal.Post) (int64, error)
	ListPost(context.Context) ([]internal.Post, error)
	GetPost(context.Context, int) (*internal.Post, error)
	DeletePost(context.Context, int) error
}


func CreatePostHandler(log *slog.Logger, queryTool PostQueryFunctions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req internal.CreatePostRequest

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
			render.JSON(w, r, utils.Error(validErr.Error()))
			return
		}
		post := internal.Post{
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




func ListPostHandler(log *slog.Logger, queryTool PostQueryFunctions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		
		posts, err := queryTool.ListPost(context.Background())
		if err != nil {
			log.Error("query error on fetching posts", logger.Err(err))
			
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, utils.Error("server error while fetchig data"))
			return
		}

		render.JSON(w, r, internal.ListPostResponse{
			Response: utils.OK(),
			Posts: posts,
		})
	}
}





func GetPostHandler(log *slog.Logger, queryTool PostQueryFunctions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		int_id, err := strconv.Atoi(id)
		if err != nil {
			log.Error("id is not integer", logger.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, utils.Error("validation error: id is not int"))
			return
		}

		post, err := queryTool.GetPost(context.Background(), int_id)
		if err != nil {
			log.Error("query error on fetching post", logger.Err(err))
			
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, utils.Error("server error while fetchig data"))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, internal.GetPostResponse{
			Response: utils.OK(),
			Post: post,
		})
	}
}



func DeletePostHandler(log *slog.Logger, queryTool PostQueryFunctions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		int_id, err := strconv.Atoi(id)
		if err != nil {
			log.Error("id is not integer", logger.Err(err))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, utils.Error("validation error: id is not int"))
			return
		}

		if err = queryTool.DeletePost(context.Background(), int_id); err != nil {
			log.Error("query error on deleting post", logger.Err(err))
			
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, utils.Error("server error while deleting post"))
			return
		} 
		
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, utils.OK())

	}
}