package redis

import (
	"context"
	"encoding/json"
	"fmt"
	goredis "github.com/redis/go-redis/v9"
	"taskManager/task/internal/model"
	"taskManager/task/pkg/redis"
	"time"
)

const (
	keyPrefix = "task:%s"
)

type Client struct {
	client *redis.Client
	ttl    time.Duration
}

func NewClient(client *redis.Client, ttl time.Duration) *Client {
	return &Client{
		client: client,
		ttl:    ttl,
	}
}

func (c *Client) Set(ctx context.Context, task model.Task) error {
	// Convert Task to TaskRedis
	taskRedis := model.TaskRedis{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	// Marshal TaskRedis instead of Task
	data, err := json.Marshal(taskRedis)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	return c.client.Unwrap().Set(ctx, c.key(task.ID), data, c.ttl).Err()
}

func (c *Client) SetMany(ctx context.Context, tasks []*model.Task) error {
	pipe := c.client.Unwrap().Pipeline()

	for _, task := range tasks {
		data, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("failed to marshal task: %w", err)
		}
		pipe.Set(ctx, c.key(task.ID), data, c.ttl)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set many tasks: %w", err)
	}
	return nil
}

func (c *Client) Get(ctx context.Context, taskID int64) (model.Task, error) {
	data, err := c.client.Unwrap().Get(ctx, c.key(taskID)).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return model.Task{}, nil // not found
		}
		return model.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	// Unmarshal to TaskRedis first
	var taskRedis model.TaskRedis
	err = json.Unmarshal(data, &taskRedis)
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	// Convert TaskRedis back to Task
	task := model.Task{
		ID:          taskRedis.ID,
		Title:       taskRedis.Title,
		Description: taskRedis.Description,
		Completed:   taskRedis.Completed,
		CreatedAt:   taskRedis.CreatedAt,
		UpdatedAt:   taskRedis.UpdatedAt,
	}

	return task, nil
}

func (c *Client) Delete(ctx context.Context, taskID int64) error {
	return c.client.Unwrap().Del(ctx, c.key(taskID)).Err()
}

func (c *Client) key(id int64) string {
	return fmt.Sprintf(keyPrefix, fmt.Sprintf("%d", id))
}
