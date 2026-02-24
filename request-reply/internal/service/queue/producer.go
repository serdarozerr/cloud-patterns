package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

const version="1.0.0"

type Producer struct{
	client* sqs.Client
	queueURL string
}

type Message struct{
	Version string `json:"version"`
	ID string `json:"id"`
	Type string `json:"type"`
	Payload map[string]any `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

func NewProducer(client *sqs.Client, queueURL string) *Producer{

	return &Producer{client: client, queueURL: queueURL}
}


func (p* Producer) SendMessage(ctx context.Context,m *Message, delaySeconds int )(string, error){
	m.Version=version
	m.Timestamp=time.Now().UTC()
	if m.ID ==""{
		m.ID=uuid.New().String()
	}

	body,err:=json.Marshal(m)
	if err !=nil{
		return "", fmt.Errorf("converting message into json: %w",err)
	}

	in:=&sqs.SendMessageInput{QueueUrl: &p.queueURL, 
		MessageBody: aws.String(string(body)), 
		DelaySeconds: int32(delaySeconds),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"MessageType":{
				DataType:aws.String("String"),
				StringValue:aws.String(m.Type),
			},
			"CorrelationId":{
				DataType:aws.String("String"),
				StringValue:aws.String(m.ID),
			},
		},
	}

	res,err:=p.client.SendMessage(ctx, in)
	if err!=nil{
		return "",fmt.Errorf("sending message: %w",err)
	}

	return *res.MessageId,nil
}


func (p* Producer)SendFIFOMessage(ctx context.Context, m *Message, messageGroupId string, deDuplicationId string) (string, error){

	m.Version=version
	m.Timestamp=time.Now().UTC()
	if m.ID ==""{
		m.ID=uuid.New().String()
	}

	body,err:=json.Marshal(m)
	if err !=nil{
		return "", fmt.Errorf("converting message into json: %w",err)
	}

	if deDuplicationId == ""{
		deDuplicationId=m.ID
	}

	in:=&sqs.SendMessageInput{
		QueueUrl: &p.queueURL,
		MessageBody: aws.String(string(body)),
		MessageGroupId: aws.String(messageGroupId),
		MessageDeduplicationId:aws.String(deDuplicationId),
	}

	result, err := p.client.SendMessage(ctx, in)
    if err != nil {
        return "", fmt.Errorf("sending FIFO message: %w", err)
    }

    return *result.MessageId, nil

}

// -------------------------------------> Batch Send

type BatchSendError struct{
	MessageID string
	Code string
	Message string
}

type BatchSendResult struct {
	Succesfull []string
	Failed []BatchSendError
}

const maxBatchSize=10

func(p* Producer) SendMessageBatch(ctx context.Context, messages []*Message)(*BatchSendResult,error){
	if len(messages) == 0{
		return nil,nil
	}

	if len(messages) >maxBatchSize{
		return nil, fmt.Errorf("batch size must be less then %d",maxBatchSize)
	}

	entries:=make([]types.SendMessageBatchRequestEntry, len(messages))

	for i, m := range messages{
		m.Version=version
		m.Timestamp=time.Now().UTC()
		if m.ID == ""{
			m.ID=uuid.New().String()
		}

		body, err := json.Marshal(m)
        if err != nil {
            return nil, fmt.Errorf("marshaling message %d: %w", i, err)
        }

		entries[i] = types.SendMessageBatchRequestEntry{
			Id: aws.String(m.ID), 
			MessageBody: aws.String(string(body)),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"MessageType":{
					DataType: aws.String("string"),
					StringValue: aws.String(m.Type),
				},
			},
		}
	}

	in:=&sqs.SendMessageBatchInput{QueueUrl: &p.queueURL, Entries: entries}

	result, err := p.client.SendMessageBatch(ctx, in)
    if err != nil {
        return nil, fmt.Errorf("batch sending messages: %w", err)
    }


	batchResult:=&BatchSendResult{
		Succesfull: make([]string, len(messages)),
		Failed: make([]BatchSendError, len(messages)),
	}

	for i, s :=range result.Successful{
		batchResult.Succesfull[i]=*s.MessageId
	}

	for i, f := range result.Failed{
		batchResult.Failed[i]=BatchSendError{
			MessageID: *f.Id,
			Code:      *f.Code,
            Message:   *f.Message,
		}
	}

	return batchResult,nil

}