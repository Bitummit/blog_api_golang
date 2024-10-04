package main

import (
	"blog_api/internal/storage/postgresql"
	"blog_api/pkg/config"
	"blog_api/pkg/logger"
	"context"
	"log/slog"
)

func main() {
	cfg := config.NewConfig()

	log := logger.NewLogger()
	log.Info("Log and config initiated", slog.Attr{Key: "env", Value: slog.StringValue(cfg.Env)})

	log.Info("Connecting database ...")
	_, err := postgresql.InitDB(context.TODO())
	if err != nil {
		log.Error("Error connection datsbase", logger.Err(err))
		return
	}
	log.Info("Success connecting database")


}