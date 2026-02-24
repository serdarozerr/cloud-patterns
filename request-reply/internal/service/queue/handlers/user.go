package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/serdarozerr/request-reply/internal/service/queue"
)
func userCreate(ctx context.Context, msg *queue.MessageConsumer)error{
	time.Sleep(100*time.Millisecond)
	slog.Info("user creation is done", "msg",msg.Payload)
	return nil
}


func userDelete(ctx context.Context, msg *queue.MessageConsumer)error{
	time.Sleep(100*time.Millisecond)
	slog.Info("user deletion is done", "msg",msg.Payload)
	return nil
}
