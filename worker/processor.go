package worker

import (
	"context"
	"github.com/hibiken/asynq"
	db "github.com/jpmoraess/pay/db/sqlc"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	store  db.Store
	server *asynq.Server
}

func NewRedisTaskProcessor(store db.Store, redisOpt asynq.RedisClientOpt) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{})
	return &RedisTaskProcessor{
		store:  store,
		server: server,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
