package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type QueueManager struct {
	client *sqs.Client
}

func NewQueuManager(client *sqs.Client) *QueueManager {

	return &QueueManager{client: client}
}

func (qm *QueueManager) CrateStandartQueue(ctx *context.Context, name string, visibilityTimeout int, messageRetantion int) (string, error) {

	in := &sqs.CreateQueueInput{
		QueueName: aws.String(name), 
		Attributes: map[string]string{
			"VisibilityTimeout": strconv.Itoa(visibilityTimeout), 
			"MessageRetationPeriot": strconv.Itoa(messageRetantion), 
			"ReceiveMessageWaitTimeSeconds": "20"}}

	out,err:=qm.client.CreateQueue(*ctx,in)
	if err != nil{
		return "", fmt.Errorf("creating standart queue %w",err)
	}

	return *out.QueueUrl, nil

}


func (qm *QueueManager) CreateFIFOQueue(ctx *context.Context, name string, visibiltyTimeout int, contentBasedDedup bool)(string, error){

	name=fmt.Sprintf("%s.fifo",name)

	attr:= map[string]string{
			"FifoQueue":"true",
			"VisibiltyTimeout":strconv.Itoa(visibiltyTimeout),
			"ReceiveMessageWaitTimeSeconds": "20",
		}
	

	if contentBasedDedup{
		attr["ContentBasedDepublication"]="true"
	}

	in:=&sqs.CreateQueueInput{
		QueueName: aws.String(name),
		Attributes: attr,}

	out,err:=qm.client.CreateQueue(*ctx,in)
	if err != nil{
		return "", fmt.Errorf("creating standart queue %w",err)
	}

	return *out.QueueUrl, nil
	
}

func (qm *QueueManager) GetQueueUrl(ctx context.Context, name string)(string ,error){

	in:=&sqs.GetQueueUrlInput{QueueName:aws.String(name)}

	res,err:=qm.client.GetQueueUrl(ctx,in)
	if err!=nil{
		return "",fmt.Errorf("getting queue url: %w",err)
	}
	return *res.QueueUrl,nil
}

func(qm *QueueManager) DeleteQueue(ctx context.Context, queueURL string) error{

	in:=&sqs.DeleteQueueInput{QueueUrl: aws.String(queueURL)}

	_, err:=qm.client.DeleteQueue(ctx,in)
	if err!=nil{
		return fmt.Errorf("deleting queue: %w",err)
	}
	return nil
}

func(qm *QueueManager) PurgeQueue(ctx context.Context, queueUrl string)error{

	in:=&sqs.PurgeQueueInput{QueueUrl: aws.String(queueUrl)}
	_,err:=qm.client.PurgeQueue(ctx,in)
	if err!=nil{
		return fmt.Errorf("purging queue: %w",err)
	}
	return nil
}

func(qm *QueueManager) ConfigureDeadLetterQueue(ctx context.Context, mainQueueUrl string, dlqARN string, maxReceiveCount int)error{

	redrivePolicy:=fmt.Sprintf( `{"deadLetterTargetArn":"%s","maxReceiveCount":"%d"}`,dlqARN,maxReceiveCount)
	in:=&sqs.SetQueueAttributesInput{QueueUrl: aws.String(mainQueueUrl),Attributes: map[string]string{"RedrivePolicy":redrivePolicy}}

	_,err:=qm.client.SetQueueAttributes(ctx, in)
	if err != nil {
        return fmt.Errorf("configuring dead letter queue: %w", err)
    }

    return nil
}

func (qm *QueueManager) GetQueueARN(ctx context.Context, queueURL string)(string,error){

	in:=&sqs.GetQueueAttributesInput{QueueUrl: aws.String(queueURL), AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn}}

	res,err:=qm.client.GetQueueAttributes(ctx,in)
	if err != nil {
        return "",fmt.Errorf("getting queue arn: %w", err)
    }

    return res.Attributes["QueueArn"], nil
}


