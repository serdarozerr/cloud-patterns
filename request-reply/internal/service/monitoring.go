package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)


type QueueStats struct{
	ApproximateNumberOfMessages int64
	ApproximateNumberOfMessagesNotVisible int64
	ApproximateNumberOfMessagesDelayed int64
}

type QueueMonitor struct{
	client *sqs.Client
}

func NewQueueMonitor(client *sqs.Client) *QueueMonitor{
	qm:= new(QueueMonitor)
	qm.client=client
	return qm
}

func(m *QueueMonitor) GetQueueStats(ctx context.Context, queueURL string)(*QueueStats, error){

	in:=&sqs.GetQueueAttributesInput{
		QueueUrl: &queueURL,
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameApproximateNumberOfMessages,
			types.QueueAttributeNameApproximateNumberOfMessagesNotVisible,
			types.QueueAttributeNameApproximateNumberOfMessagesDelayed,
		},
	}

	res,err:= m.client.GetQueueAttributes(ctx,in)
	if err !=nil{
		return nil,fmt.Errorf("getting queueu attrb: %w",err)
	}

	stats :=&QueueStats{}
	if val,ok:=res.Attributes["ApproximateNumberOfMessages"]; ok{
		stats.ApproximateNumberOfMessages,_=strconv.ParseInt(val,10,64)
	}

	if val,ok:=res.Attributes["ApproximateNumberOfMessagesNotVisible"]; ok{
		stats.ApproximateNumberOfMessagesNotVisible,_=strconv.ParseInt(val,10,64)
	}

	if val,ok:=res.Attributes["ApproximateNumberOfMessagesDelayed"]; ok{
		stats.ApproximateNumberOfMessagesDelayed,_=strconv.ParseInt(val,10,64)
	}

	return stats,nil
}

func (m *QueueMonitor) HealthCheck(ctx context.Context, queueURL string) error{
	_, err := m.GetQueueStats(ctx, queueURL)
    if err != nil {
        return fmt.Errorf("queue health check failed: %w", err)
    }
    return nil
}