package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	example "github.com/yybirdcf/micro/example/proto/example"
)

type Example struct{}

func (e *Example) Handle(ctx context.Context, msg *example.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *example.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
