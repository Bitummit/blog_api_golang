package api

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Bitummit/blog_api_golang/internal"
	"github.com/Bitummit/blog_api_golang/pkg/config"
	"github.com/Bitummit/blog_api_golang/pkg/logger"
	"github.com/Bitummit/blog_api_golang/pkg/utils"
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


type HTTPServer struct {
	Log *slog.Logger
	Storage PostQueryFunctions
	Cfg *config.Config
	Router chi.Router
}


func StartServer(server *HTTPServer) error{

	server.Router.Use(middleware.RequestID)
	server.Router.Use(middleware.RealIP)
	server.Router.Use(middleware.Logger)
	server.Router.Use(middleware.Recoverer)
	server.Router.Use(middleware.URLFormat)
	server.Router.Use(internal.SetJSONContentType)

	server.Router.Post("/post/", server.CreatePostHandler)
	server.Router.With(
		internal.CheckTokenMiddleware(server.Log),
	).Get("/post/", server.ListPostHandler,)
	server.Router.Get("/post/{id}/", server.GetPostHandler)
	server.Router.Delete("/post/{id}/", server.DeletePostHandler)


	srv := &http.Server{
		Addr: server.Cfg.Address,
		Handler: server.Router,
		ReadTimeout: server.Cfg.Timeout,
		WriteTimeout: server.Cfg.Timeout,

	}
	server.Log.Info("Starting server ...", slog.String("address", server.Cfg.Address))

	if err := srv.ListenAndServe(); err != nil {
		server.Log.Error("Can't start server", logger.Err(err))
		return err
	}

	// Graceful shutdown
	server.Log.Info("Server stopped")
	return nil
}


func (s *HTTPServer) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
		s.Log = slog.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req internal.CreatePostRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			s.Log.Error("failed to decode request", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, utils.Error("error decoding request"))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validErr := err.(validator.ValidationErrors)

			s.Log.Error("Validate error", logger.Err(validErr))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, utils.Error(validErr.Error()))
			return
		}
		post := internal.Post{
			Title: req.Title,
			Body: req.Body,
			Author: req.Author,
		}

		id, err := s.Storage.NewPost(context.Background(), post)
		if err != nil {
			s.Log.Error("Error while adding new post", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, utils.Error("db error"))
			return
		}
		s.Log.Info("Inserted post", slog.Int64("id", int64(id)))

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, utils.OK())
}


func (s *HTTPServer) ListPostHandler(w http.ResponseWriter, r *http.Request) {
	s.Log = slog.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	
	posts, err := s.Storage.ListPost(context.Background())
	if err != nil {
		s.Log.Error("query error on fetching posts", logger.Err(err))
		
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, utils.Error("server error while fetchig data"))
		return
	}

	render.JSON(w, r, internal.ListPostResponse{
		Response: utils.OK(),
		Posts: posts,
	})
}


func (s *HTTPServer) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	s.Log = slog.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	id := chi.URLParam(r, "id")
	int_id, err := strconv.Atoi(id)
	if err != nil {
		s.Log.Error("id is not integer", logger.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, utils.Error("validation error: id is not int"))
		return
	}

	post, err := s.Storage.GetPost(context.Background(), int_id)
	if err != nil {
		s.Log.Error("query error on fetching post", logger.Err(err))
		
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



func (s *HTTPServer) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	s.Log = slog.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	id := chi.URLParam(r, "id")
	int_id, err := strconv.Atoi(id)
	if err != nil {
		s.Log.Error("id is not integer", logger.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, utils.Error("validation error: id is not int"))
		return
	}

	if err = s.Storage.DeletePost(context.Background(), int_id); err != nil {
		s.Log.Error("query error on deleting post", logger.Err(err))
		
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, utils.Error("server error while deleting post"))
		return
	} 
	
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, utils.OK())
}
