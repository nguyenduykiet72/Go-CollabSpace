package worker

import (
	"Go-CollabSpace/pkg/logger"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskUpdateSearchIndex(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	db     *gorm.DB
	// TODO - add AI service OpenAI/local models
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, db *gorm.DB) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{Concurrency: 10},
	)

	return &RedisTaskProcessor{server: server, db: db}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux() // Register task handlers
	mux.HandleFunc(TaskUpdateSearchIndex, processor.ProcessTaskUpdateSearchIndex)
	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) ProcessTaskUpdateSearchIndex(ctx context.Context, task *asynq.Task) error {
	var payload PayloadUpdateSearchIndex
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("Processing TaskUpdateSearchIndex for DocID: %v in WorkspaceID: %v", payload.DocID, payload.WorkspaceID))

	return nil
}
