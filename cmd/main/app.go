package main

import (
	"context"
	"log/slog"
	// "os"
	// "os/signal"
	// "syscall"

	"github.com/Bitummit/blog_api_golang/internal/api"
	// blogservice "github.com/Bitummit/blog_api_golang/internal/blog_service"
	"github.com/Bitummit/blog_api_golang/internal/storage/postgresql"
	// "github.com/IBM/sarama"

	_ "github.com/Bitummit/blog_api_golang/docs"
	"github.com/Bitummit/blog_api_golang/pkg/config"
	"github.com/Bitummit/blog_api_golang/pkg/logger"
	"github.com/go-chi/chi/v5"
)


// var sigchan = make(chan os.Signal, 1)
// var doneCh = make(chan struct{})


//	@title			Go Blog API
//	@version		1.0
//	@description	This is a sample API blog service.

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						JWT

//	@host		localhost:8000
//	@BasePath	/
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

	// <-doneCh
	
	// if err := worker.Close(); err != nil {
	// 	panic(err)
	// } 
	

}


// func runWorker(log *slog.Logger, consumer sarama.PartitionConsumer) {
// 	for {
// 		select {
// 		case err := <- consumer.Errors():
// 			log.Error("Consumer error", err)
// 		case msg := <-consumer.Messages():
// 			post := string(msg.Value)
// 			log.Info("got new post!", post)
// 		case <- sigchan:
// 			doneCh <- struct{}{}
// 		}
// 	}
// }
// TODO: filtering
// TODO: sorting
// TODO: pagination
