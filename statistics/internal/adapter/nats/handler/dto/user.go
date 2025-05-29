package dto

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"taskManager/proto/gen/statistics"
	"taskManager/statistics/internal/model"
)

func ToUser(msg *nats.Msg) (model.UserNats, error) {
	var pbClient statistics.GetUserStatisticsResponse
	err := proto.Unmarshal(msg.Data, &pbClient)
	if err != nil {
		return model.UserNats{}, fmt.Errorf("proto.Unmarshall: %w", err)
	}

	return model.UserNats{
		ID: pbClient.Id,
	}, nil
}
