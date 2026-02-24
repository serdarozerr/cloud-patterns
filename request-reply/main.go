package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/serdarozerr/request-reply/internal/api"
	"github.com/serdarozerr/request-reply/internal/config"
	"github.com/serdarozerr/request-reply/internal/service/queue"
	"github.com/serdarozerr/request-reply/internal/service/queue/handlers"
)

func configureLogger() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))

}

func getSqsClient(awsCfg *config.AWSConfig) *sqs.Client{

	ctx:=context.Context(context.Background())

	// create queue client
	client,err:=queue.NewSQSClient(ctx,awsCfg)
	if err !=nil{
		slog.Error("Failed to create queue client","error",err)
		panic(1)
	}
	return client
}
func getQueueURL(ctx context.Context, client *sqs.Client, awsCfg *config.AWSConfig) string{
	queueMgr:=queue.NewQueuManager(client)

	queueUrl:=awsCfg.QueueURL
	var err error
	if awsCfg.QueueURL ==""{
		queueUrl,err=queueMgr.GetQueueUrl(ctx,awsCfg.Name)
		if err!=nil{
			queueUrl,err=queueMgr.CrateStandartQueue(ctx,awsCfg.Name, 30, 345600)
			if err!=nil{
				slog.Error("Failed to create queue","error",err)
				panic(1)
			}
			slog.Info("Qeueu created")
		}
	}
	return queueUrl
}

func getProducerQueue(ctx context.Context,awsCfg *config.AWSConfig) *queue.Producer{
	client:=getSqsClient(awsCfg)
	queueUrl:=getQueueURL(ctx,client,awsCfg)
	awsCfg.QueueURL=queueUrl
	prod:=queue.NewProducer(client,queueUrl)
	return prod
}




func getConsumerQueue(ctx context.Context, awsCfg *config.AWSConfig) *queue.Consumer{
	client:=getSqsClient(awsCfg)
	queueUrl:=getQueueURL(ctx, client, awsCfg)
	cons:=queue.NewConsumer(client,
	queue.ConsumerConfig{
		QueueURL:          queueUrl,
        MaxMessages:       10,
        VisibilityTimeout: 30,
        WaitTimeSeconds:   20,
        WorkerCount:       5,
	},
	handlers.MessageHandler)
	return cons
}

// This mode start a server with endpoints that
// time taking tasks/jobs will be passed to queue
// to send worker server
func startProducerServer(cfg *config.Config, awsCfg *config.AWSConfig){
	slog.Info("Starting server", "host", cfg.Host, "port",cfg.Port)

	producer:=getProducerQueue(context.Background(),awsCfg)

	s := http.Server{
		Addr:    fmt.Sprintf("%s:%s",cfg.Host,cfg.Port),
		Handler: api.NewRouter(producer),
		ReadTimeout: 10 *time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func ()  {
		if err := s.ListenAndServe(); err != nil {
			slog.Error("Failed to start server", "port", 8080)
		}
	}()

	// create os.Signal type channel,
	// send signal to chanell when term or int
	quit :=make(chan os.Signal,1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	slog.Info("Shutting down producer")

	shutdownContext, shutdownCancel :=context.WithTimeout(context.Background(), 30 *time.Second)
	defer shutdownCancel()

	if err:=s.Shutdown(shutdownContext); err!=nil{
		slog.Info("Server shutdown error", "error",err)
	}

	slog.Info("Shutdown completed")
}


// This mode is worker mode, listens the queue
// and process any task/job, no endpoints exposes
// in this mode
func startConsumerWorker(awsCfg *config.AWSConfig){
	ctx:=context.Background()
	consumer:=getConsumerQueue(ctx,awsCfg)
	go func ()  {
		slog.Info("starting consumer")
		if err:=consumer.Start(ctx); err != nil{
			slog.Info("error while starting consumer", "error",err)
		}
	}()

	// create os.Signal type channel,
	// send signal to channel when term or int
	quit :=make(chan os.Signal,1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	slog.Info("Shutting down consumer")

}

func main() {
	configureLogger()
	cfg:=config.NewConfig("config.json")
	awsCfg:=config.NewAwsConfig("aws_cred.json")

	switch cfg.Mode {
	case "producer":
		startProducerServer(cfg,awsCfg)
	case "consumer":
		startConsumerWorker(awsCfg)
	default:
		slog.Info("Unsupported mode, supported modes are: producer, consumer")
		panic(1)	
	}
	

}
