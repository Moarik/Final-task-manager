package dto

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"taskManager/proto/gen/statistics"
	"taskManager/statistics/internal/model"
)

func ToTask(msg *nats.Msg) (model.TaskNats, error) {
	var pbClient statistics.GetTaskStatisticsResponse
	err := proto.Unmarshal(msg.Data, &pbClient)
	if err != nil {
		return model.TaskNats{}, fmt.Errorf("proto.Unmarshall: %w", err)
	}

	return model.TaskNats{
		ID: pbClient.Id,
	}, nil
}
