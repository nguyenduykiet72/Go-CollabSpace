package worker

import (
	"Go-CollabSpace/internal/common/infrastructure"
	"Go-CollabSpace/pkg/logger"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskUpdateSearchIndex(ctx context.Context, task *asynq.Task) error
	ProcessTaskSendResetEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	db     *gorm.DB
	// TODO - add AI service OpenAI/local models
	emailSender infrastructure.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, db *gorm.DB, emailSender infrastructure.EmailSender) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{Concurrency: 10},
	)

	return &RedisTaskProcessor{
		server:      server,
		db:          db,
		emailSender: emailSender,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux() // Register task handlers
	mux.HandleFunc(TaskUpdateSearchIndex, processor.ProcessTaskUpdateSearchIndex)
	mux.HandleFunc(TaskSendResetEmail, processor.ProcessTaskSendResetEmail)
	return processor.server.Start(mux)
}

func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}

func (processor *RedisTaskProcessor) ProcessTaskUpdateSearchIndex(ctx context.Context, task *asynq.Task) error {
	var payload PayloadUpdateSearchIndex
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("Processing TaskUpdateSearchIndex for DocID: %v in WorkspaceID: %v", payload.DocID, payload.WorkspaceID))

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendResetEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendResetEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	err := processor.emailSender.SendResetPasswordEmail(payload.ToEmail, payload.ResetToken)
	if err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("Sent password reset email to: %v", payload.ToEmail))
	return nil
}
