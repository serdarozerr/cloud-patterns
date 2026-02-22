package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/serdarozerr/request-reply/internal/api"
	"github.com/serdarozerr/request-reply/internal/config"
	"github.com/serdarozerr/request-reply/internal/service"
)

func configureLogger() {

	handler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(handler))

}

func getSqsClient(awsCfg *config.AWSConfig) *sqs.Client{

	ctx:=context.Context(context.Background())

	// create queue client
	client,err:=service.NewSQSClient(ctx,awsCfg)
	if err !=nil{
		slog.Error("Failed to create queue client: %w",err)
		panic(1)
	}
	return client
}
func getQueueURL(ctx context.Context, client *sqs.Client, awsCfg *config.AWSConfig) string{

	// create queue manager
	queueMgr:=service.NewQueuManager(client)

	queueUrl:=awsCfg.QueueURL
	var err error
	if awsCfg.QueueURL ==""{
		queueUrl,err=queueMgr.GetQueueUrl(ctx,awsCfg.Name)
		if err!=nil{
			queueUrl,err=queueMgr.CrateStandartQueue(ctx,awsCfg.Name, 30, 345600)
			if err!=nil{
				slog.Error("Failed to create queue: %w",err)
				panic(1)
			}
			slog.Info("Qeueu created")
		}
	}
	return queueUrl
}

func getProducerQueue(ctx context.Context,awsCfg *config.AWSConfig) *service.Producer{
	client:=getSqsClient(awsCfg)
	queueUrl:=getQueueURL(ctx,client,awsCfg)
	awsCfg.QueueURL=queueUrl
	prod:=service.NewProducer(client,queueUrl)
	return prod
}


func userRegister(ctx context.Context, msg *service.MessageConsumer)error{
	time.Sleep(100*time.Millisecond)
	slog.Info("user registeration is done", "msg",msg.Payload)
	return nil
}


func userDelete(ctx context.Context, msg *service.MessageConsumer)error{
	time.Sleep(100*time.Millisecond)
	slog.Info("user deletion is done", "msg",msg.Payload)
	return nil
}

func messageHandler(ctx context.Context, msg *service.MessageConsumer)error{
	slog.Info("Processing message", "id", msg.ID, "type", msg.Type)

	switch msg.Type{
	case "user.register":
		return userRegister(ctx,msg)
	case "user.delete":
		return userDelete(ctx,msg)
	default:
		slog.Info("Unknown message type" , "type", msg.Type)
		return nil
	}
}

func getConsumerQueue(ctx context.Context, awsCfg *config.AWSConfig) *service.Consumer{
	client:=getSqsClient(awsCfg)
	queueUrl:=getQueueURL(ctx, client, awsCfg)
	cons:=service.NewConsumer(client,
	service.ConsumerConfig{
		QueueURL:          queueUrl,
        MaxMessages:       10,
        VisibilityTimeout: 30,
        WaitTimeSeconds:   20,
        WorkerCount:       5,
	},
	messageHandler)
	return cons
}

// This mode start a server with endpoints that
// time taking tasks/jobs will be passed to queue
// to send worker server
func startProducerServer(){
	slog.Info("Starting server on port 8080")
	s := http.Server{
		Addr:    ":8080",
		Handler: api.NewRouter(),
	}

	if err := s.ListenAndServe(); err != nil {
		slog.Error("Failed to start server", "port", 8080)
	}
}


// This mode is worker mode, listens the queue
// and process any task/job, no endpoints exposes
// in this mode
func startConsumerWorker(){

}

func main() {
	configureLogger()
	

}
