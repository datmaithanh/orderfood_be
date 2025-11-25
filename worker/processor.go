package worker

import (
	"context"

	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/hibiken/asynq"
)

type TaskProcessor interface {
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task ) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
		},
	)
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
