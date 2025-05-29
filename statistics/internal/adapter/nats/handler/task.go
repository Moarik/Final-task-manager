package handler

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"taskManager/statistics/internal/adapter/nats/handler/dto"
)

type Client struct {
	statisticsUC StatisticsUsecase
}

func NewClient(statisticsUC StatisticsUsecase) *Client {
	return &Client{
		statisticsUC: statisticsUC,
	}
}

func (c *Client) HandleUserCreated(ctx context.Context, msg *nats.Msg) error {
	user, err := dto.ToUser(msg)
	if err != nil {
		log.Println("failed to decode user message:", err)
		return err
	}

	if err = c.statisticsUC.CreateUserStatistics(ctx, user); err != nil {
		log.Println("failed to create user:", err)
		return err
	}

	return nil
}

func (c *Client) HandleTaskCreated(ctx context.Context, msg *nats.Msg) error {
	task, err := dto.ToTask(msg)
	if err != nil {
		log.Println("failed to decode task message:", err)
		return err
	}

	fmt.Println("FROM TASK NATS:", task)

	if err = c.statisticsUC.CreateTaskStatistics(ctx, task); err != nil {
		log.Println("failed to create task:", err)
		return err
	}

	return nil
}
