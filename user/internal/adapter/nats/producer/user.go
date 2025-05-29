package producer

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"taskManager/user/internal/adapter/nats/producer/dto"
	"taskManager/user/internal/model"
	"taskManager/user/pkg/nats"
	"time"
)

const PushTimeout = time.Second * 30

type Client struct {
	client  *nats.Client
	subject string
}

func NewUserProducer(
	client *nats.Client,
	subject string,
) *Client {
	return &Client{
		client:  client,
		subject: subject,
	}
}

func (c *Client) Push(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, PushTimeout)
	defer cancel()

	// TODO after proto file
	clientPb := dto.FromClient(user)
	data, err := proto.Marshal(&clientPb)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %w", err)
	}

	err = c.client.Conn.Publish(c.subject, data)
	if err != nil {
		return fmt.Errorf("c.client.Conn.Publish: %w", err)
	}
	log.Println("client is pushed:", user)

	return nil
}
