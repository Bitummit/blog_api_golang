package main

import (
	"context"
	"log/slog"

	"github.com/Bitummit/blog_api_golang/internal/api"
	"github.com/Bitummit/blog_api_golang/internal/storage/postgresql"
	"github.com/Bitummit/blog_api_golang/pkg/config"
	"github.com/Bitummit/blog_api_golang/pkg/logger"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.NewConfig()

	log := logger.NewLogger()
	log.Info("Log and config initiated", slog.Attr{Key: "env", Value: slog.StringValue(cfg.Env)})

	log.Info("Connecting database ...")
	storage, err := postgresql.InitDB(context.TODO())
	if err != nil {
		log.Error("Error connection database", logger.Err(err))
		return
	}
	defer storage.DB.Close()

	log.Info("Success connecting database")

	router := chi.NewRouter()

	server := api.HTTPServer{
		Log: log,
		Storage: storage,
		Cfg: cfg,
		Router: router,

	}
	if err := api.StartServer(&server); err != nil {
		log.Error("Server error", logger.Err(err))
		// Graceful shutdown
	}

}

// TODO: filtering
// TODO: sorting
// TODO: pagination
