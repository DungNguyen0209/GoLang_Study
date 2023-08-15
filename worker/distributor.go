package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

// return to force the type return must implement interface
type TaskDistributor interface {
	DistributeTaskVerifyEmail(
		ctx context.Context,
		payload *PayLoadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	// Send Task to send task to Redis queue
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
