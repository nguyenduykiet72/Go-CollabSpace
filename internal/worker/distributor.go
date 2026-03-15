package worker

import (
	"Go-CollabSpace/pkg/logger"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const TaskUpdateSearchIndex = "task:update_search_index"
const TaskSendResetEmail = "task:send_reset_email"

type PayloadUpdateSearchIndex struct {
	DocID       uuid.UUID `json:"doc_id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	PlainText   string    `json:"plain_text"`
}

type PayloadSendResetEmail struct {
	ToEmail    string `json:"to_email"`
	ResetToken string `json:"reset_token"`
}

type TaskDistributor interface {
	DistributeTaskUpdateSearchIndex(ctx context.Context, payload *PayloadUpdateSearchIndex, opts ...asynq.Option) error
	DistributeTaskSendResetEmail(ctx context.Context, payload *PayloadSendResetEmail, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{client: client}
}

func (d *RedisTaskDistributor) DistributeTaskUpdateSearchIndex(ctx context.Context, payload *PayloadUpdateSearchIndex, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	task := asynq.NewTask(TaskUpdateSearchIndex, jsonPayload, opts...)

	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("Enqueued task: %v", info))
	return nil
}

func (d *RedisTaskDistributor) DistributeTaskSendResetEmail(ctx context.Context, payload *PayloadSendResetEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	task := asynq.NewTask(TaskSendResetEmail, jsonPayload, opts...)

	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	logger.Log.Info(fmt.Sprintf("Enqueued task: %v", info))
	return nil
}
