package main

import run "github.com/Bitummit/blog_api_golang/internal"

//	@title			Go Blog API
//	@version		1.0
//	@description	This is a sample API blog service.

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						JWT

//	@host		localhost:8000
//	@BasePath	/
func main() {
	run.Run()
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
