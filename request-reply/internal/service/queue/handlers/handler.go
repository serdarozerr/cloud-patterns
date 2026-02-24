package handlers

import (
	"context"
	"log/slog"

	"github.com/serdarozerr/request-reply/internal/service/queue"
)

func MessageHandler(ctx context.Context, msg *queue.MessageConsumer)error{
	slog.Info("Processing message", "id", msg.ID, "type", msg.Type)

	switch msg.Type{
	case "user.create":
		return userCreate(ctx,msg)
	case "user.delete":
		return userDelete(ctx,msg)
	default:
		slog.Info("Unknown message type" , "type", msg.Type)
		return nil
	}
}