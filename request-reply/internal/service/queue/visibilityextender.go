package queue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type VisibilityExtender struct{
	client *sqs.Client
	queueURL string
}

func NewVisibilityExtender( client *sqs.Client, queueURL string) *VisibilityExtender{

	return &VisibilityExtender{
		client: client,
		queueURL: queueURL,
	}
}

func (v *VisibilityExtender) ExtendVisibility(ctx context.Context, receiptHandle string, additionalSeconds int) error{


	in:=&sqs.ChangeMessageVisibilityInput{QueueUrl: &v.queueURL, ReceiptHandle: &receiptHandle, VisibilityTimeout: int32(additionalSeconds)}

	_, err:=v.client.ChangeMessageVisibility(ctx, in)

	if err != nil{
		return fmt.Errorf("extending visibility timeout: %w", err)
	}

	return nil
}

func (v *VisibilityExtender) StartVisibilityHeartbeat(ctx context.Context, receiptHandle string, interval time.Duration, extension int){

	ticker := time.NewTicker(interval)
    defer ticker.Stop()

	for {
		select{
			case <-ctx.Done():
            	return
			case <-ticker.C:
				 if err := v.ExtendVisibility(ctx, receiptHandle, extension); err != nil {
                	log.Printf("Failed to extend visibility: %v", err)
                	return
            	}
		}
	}
}