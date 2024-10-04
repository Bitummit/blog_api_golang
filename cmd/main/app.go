package main

import (
	"blog_api/internal/post"
	"blog_api/internal/storage/postgresql"
	"blog_api/pkg/config"
	"blog_api/pkg/logger"
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.NewConfig()

	log := logger.NewLogger()
	log.Info("Log and config initiated", slog.Attr{Key: "env", Value: slog.StringValue(cfg.Env)})

	log.Info("Connecting database ...")
	storage, err := postgresql.InitDB(context.TODO())
	if err != nil {
		log.Error("Error connection datsbase", logger.Err(err))
		return
	}
	log.Info("Success connecting database")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/post", post.CreatePost(log, storage))

	srv := &http.Server{
		Addr: cfg.Address,
		Handler: router,
		ReadTimeout: cfg.Timeout,
		WriteTimeout: cfg.Timeout,

	}
	log.Info("Starting server ...", slog.String("address", cfg.Address))

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Cant start server", logger.Err(err))
	}

	log.Info("Server stopped")

}