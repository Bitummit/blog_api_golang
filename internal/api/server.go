package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Bitummit/blog_api_golang/internal"
	blogservice "github.com/Bitummit/blog_api_golang/internal/blog_service"
	"github.com/Bitummit/blog_api_golang/internal/models"

	"github.com/Bitummit/blog_api_golang/pkg/config"
	"github.com/Bitummit/blog_api_golang/pkg/logger"
	"github.com/Bitummit/blog_api_golang/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type PostService interface {
	CreatePostService(blogservice.PostQueryFunctions, models.Post) (int64, error)
	ListPostService(blogservice.PostQueryFunctions)  ([]models.Post, error)
	GetPostService(blogservice.PostQueryFunctions, int) (*models.Post, error)
	DeletePostService(blogservice.PostQueryFunctions, int) error
}


type HTTPServer struct {
	Log *slog.Logger
	Storage blogservice.PostQueryFunctions
	Cfg *config.Config
	Router chi.Router
}


func StartServer(server *HTTPServer) error{

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx) // TODO: make it useful
	defer cancel()

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

	server.Router.Post("/login/", server.LoginHandler)

	srv := &http.Server{
		Addr: server.Cfg.Address,
		Handler: server.Router,
		ReadTimeout: server.Cfg.Timeout,
		WriteTimeout: server.Cfg.Timeout,

	}
	server.Log.Info("Starting server ...", slog.String("address", server.Cfg.Address))

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server stopped %v", err)
	}

	// Graceful shutdown
	server.Log.Info("Server stopped")
	return nil
}


func (s *HTTPServer) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	s.Log = slog.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req CreatePostRequest
	err := render.DecodeJSON(r.Body, &req)
	r.Body.Close()

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

	post := models.Post{
		Title: req.Title,
		Body: req.Body,
		Author: req.Author,
	}
	id, err := blogservice.CreatePostService(s.Storage, post)
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

	posts, err := blogservice.ListPostService(s.Storage)
	if err != nil {
		s.Log.Error("query error on fetching posts", logger.Err(err))
		
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, utils.Error("server error while fetchig data"))
		return
	}

	render.JSON(w, r, ListPostResponse{
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

	post, err := blogservice.GetPostService(s.Storage, int_id)
	if err != nil {
		s.Log.Error("query error on fetching post", logger.Err(err))
		
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, utils.Error("no such post"))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, GetPostResponse{
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

	if err = blogservice.DeletePostService(s.Storage, int_id); err != nil {
		s.Log.Error("query error on deleting post", logger.Err(err))
		
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, utils.Error("server error while deleting post"))
		return
	} 
	
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, utils.OK())
}


func (s *HTTPServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	s.Log = slog.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req LoginRequest
	err := render.DecodeJSON(r.Body, &req)
	r.Body.Close()
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

	user := models.User{
		Username: req.Username,
		Password: req.Password,
	}
	token, err := blogservice.LoginService(s.Storage, s.Log, user)
	if err != nil {

		// *handle different error types!*

		s.Log.Error("Error while loggining", logger.Err(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, utils.Error("server error"))
		return
	}

	resp := LoginResponse{
		Response: utils.OK(),
		Token: *token,
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp)

}