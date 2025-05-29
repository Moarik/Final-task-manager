package producer

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"taskManager/task/internal/adapter/nats/producer/dto"
	"taskManager/task/internal/model"
	"taskManager/task/pkg/nats"
	"time"
)

const PushTimeout = time.Second * 30

type Client struct {
	client  *nats.Client
	subject string
}

func NewTaskProducer(
	client *nats.Client,
	subject string,
) *Client {
	return &Client{
		client:  client,
		subject: subject,
	}
}

func (c *Client) Push(ctx context.Context, task model.Task) error {
	ctx, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	clientPb := dto.FromTask(task)
	data, err := proto.Marshal(&clientPb)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	err = c.client.Conn.Publish(c.subject, data)
	if err != nil {
		return fmt.Errorf("c.client.Conn.Publish: %w", err)
	}
	log.Println("task is pushed:", task)

	return nil
}
