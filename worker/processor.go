package worker

import (
	"context"

	db "github.com/arya2004/xyfin/db/sqlc"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	DefaultQueue = "default"	
)

type TaskProcessor interface {
	Start()	error
	ProcessVerifyEmail(ctx context.Context, task *asynq.Task) error
	
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor{
	server := asynq.NewServer(
		redisOpt, 
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				DefaultQueue: 5,
			},
		},
	)
	return &RedisTaskProcessor{server: server, store: store}
}


func (processor *RedisTaskProcessor) Start() error {

	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessVerifyEmail)

	return processor.server.Start(mux)
}