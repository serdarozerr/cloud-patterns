package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type MessageConsumer struct{
	ID string `json:"id"`
	Type string `json:"type"`
	Payload map[string]any `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
	ReceiptHandle string `json:"-"`
	Attributes map[string]string `json:"-"`
}

type Handler func(ctx context.Context, m* MessageConsumer) error

type Consumer struct{
	client* sqs.Client
	queueURL string
	handler Handler
	// total number of messages 
	// in one call
	maxMessages int
	// message become visible again
	// after this timeout, so other
	// workers can see the task again
	visibilityTimeout int
	// polling timeout,if there is no
	// message wait this many sec.
	// if message arrives 
	// wake up immediately
	waitTimeSeconds int
	workerCount int
}

type ConsumerConfig struct{
	QueueURL string
	MaxMessages int
	VisibilityTimeout int
	WaitTimeSeconds int
	WorkerCount int
}

func NewConsumer(client *sqs.Client, cfg ConsumerConfig, handler Handler) *Consumer{

	if cfg.MaxMessages <=0 || cfg.MaxMessages >10{
		cfg.MaxMessages=10
	}

	if cfg.VisibilityTimeout<=0{
		cfg.VisibilityTimeout=30
	}

	if cfg.WaitTimeSeconds<=0{
		cfg.WaitTimeSeconds=20
	}

	if cfg.WorkerCount <=0{
		cfg.WorkerCount=5
	}

	return &Consumer{
		client: client,
		handler: handler,
		queueURL: cfg.QueueURL,
		maxMessages: cfg.MaxMessages,
		waitTimeSeconds: cfg.WaitTimeSeconds,
		workerCount: cfg.WorkerCount,
	}

}

func (c *Consumer) Start(ctx context.Context)error{

	wg:=sync.WaitGroup{}
	msgChan := make(chan *MessageConsumer, c.workerCount*2)
	for i :=range c.workerCount{
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			c.work(ctx,msgChan)

		}(i)
	}

	go func(){
		c.pool(ctx,msgChan)
		close(msgChan)
	}()

	wg.Wait()

	return ctx.Err()
}

func (c *Consumer) work(ctx context.Context, msgChan <-chan *MessageConsumer){

	for msg := range msgChan{
		select{
		case <-ctx.Done():
			return
		default:
			c.processMessage(ctx, msg)
		}
	}
}

func (c *Consumer) pool(ctx context.Context, msgChan chan<- *MessageConsumer){
	for  {
		select {
			case <-ctx.Done():
				return
			default:
				messages,err :=c.receiveMessages(ctx)
				if err != nil{
					slog.Error("receiving messages","error",err)
					time.Sleep(time.Second)
					continue
				}

				for _,m := range messages{
					select{
					case msgChan <-m:
					case <-ctx.Done():
						return
					}
				}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg *MessageConsumer){

	ctxT,cancel:=context.WithTimeout(ctx,time.Duration(time.Duration(c.visibilityTimeout-5)*time.Second))
	defer cancel()

	err:=c.handler(ctxT,msg)
	 if err != nil {
        slog.Info("Error processing message %s: %v", msg.ID, err)
        return
    }

	 if err := c.deleteMessage(ctx, msg.ReceiptHandle); err != nil {
        slog.Info("Error deleting message %s: %v", msg.ID, err)
    }
}


func (c *Consumer) receiveMessages(ctx context.Context)([]*MessageConsumer, error){

in:=&sqs.ReceiveMessageInput{
	QueueUrl: &c.queueURL,
	MaxNumberOfMessages: int32(c.maxMessages),
	VisibilityTimeout: int32(c.visibilityTimeout),
	WaitTimeSeconds: int32(c.waitTimeSeconds),
	MessageAttributeNames: []string{"All"},
	AttributeNames: []types.QueueAttributeName{
            types.QueueAttributeNameAll,
        },
}

res,err:=c.client.ReceiveMessage(ctx,in)
if err != nil{
	return nil, fmt.Errorf("receiving messages: %w", err)
}

messages :=make([]*MessageConsumer, len(res.Messages))
for i, m :=range res.Messages{
	 var msg MessageConsumer
    if err := json.Unmarshal([]byte(*m.Body), &msg); err != nil {
        log.Printf("Error unmarshaling message %s: %v", *m.MessageId, err)
        continue
    }

    msg.ReceiptHandle = *m.ReceiptHandle
    msg.Attributes = m.Attributes
    messages[i] = &msg
}
return messages, nil
}

func (c *Consumer) deleteMessage(ctx context.Context, rh string )error{

	in:=&sqs.DeleteMessageInput{
		QueueUrl: &c.queueURL,
		ReceiptHandle: &rh,
	}

	_,err:=c.client.DeleteMessage(ctx,in)

	if err!=nil{
		return fmt.Errorf("deleting message: %w",err)
	}
	return nil
}

//-----------------------------> Batch Delete

type BatchDeleteError struct{
	ID string
	Code string
	Message string
}

type BatchDeleteResult struct{
	Successful []string
	Failed []BatchDeleteError
}

func (c *Consumer) DeleteMessageBatch(ctx context.Context, messages []*MessageConsumer)(*BatchDeleteResult, error){

	if len(messages) ==0{
		return &BatchDeleteResult{}, nil
	}
	if len(messages) >10{
		return nil, fmt.Errorf("batch size exceeds maximum of 10 messages")
	}

	entries := make([]types.DeleteMessageBatchRequestEntry, len(messages))
	for i, m := range messages{
		entries[i] = types.DeleteMessageBatchRequestEntry{
			Id: aws.String(m.ID),
			ReceiptHandle: aws.String(m.ReceiptHandle),
		}
	}

	in:=&sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(c.queueURL),
		Entries: entries,
	}
	  
	result, err := c.client.DeleteMessageBatch(ctx, in)
    if err != nil {
        return nil, fmt.Errorf("batch deleting messages: %w", err)
    }

	batchResult := &BatchDeleteResult{
        Successful: make([]string, len(result.Successful)),
        Failed:     make([]BatchDeleteError, len(result.Failed)),
    }

    for i, s := range result.Successful {
        batchResult.Successful[i] = *s.Id
    }

    for i, f := range result.Failed {
        batchResult.Failed[i] = BatchDeleteError{
            ID:      *f.Id,
            Code:    *f.Code,
            Message: *f.Message,
        }
    }

    return batchResult, nil


}