package run

import (
	"context"
	"log/slog"

	"github.com/Bitummit/blog_api_golang/internal/api/server"
	"github.com/Bitummit/blog_api_golang/internal/storage/postgresql"

	_ "github.com/Bitummit/blog_api_golang/docs"
	"github.com/Bitummit/blog_api_golang/pkg/config"
	"github.com/Bitummit/blog_api_golang/pkg/logger"
	"github.com/go-chi/chi/v5"
)


func Run() {
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
	srv := server.HTTPServer{
		Log: log,
		Storage: storage,
		Cfg: cfg,
		Router: router,

	}
	if err := server.StartServer(&srv); err != nil {
		log.Error("Server error", logger.Err(err))
	}
}

// kafka := blogservice.NewKafka(log)
	// topic := "new_posts"
	// worker, err := kafka.ConnectConsumer([]string{"localhost:9092"})
	// if err != nil {
	// 	log.Error("Error connecting to kafka consumer", err)
	// 	return
	// }

	// consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	// if err != nil {
	// 	log.Error("Error connecting to kafka consumer", err)
	// 	return
	// }
	// signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT)
	// go runWorker(log, consumer)