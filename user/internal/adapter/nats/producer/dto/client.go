package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"taskManager/proto/gen/user"
	"taskManager/user/internal/model"
)

func FromClient(client model.User) user.UserCreateNats {
	return user.UserCreateNats{
		Id:        client.ID,
		Name:      client.Name,
		Email:     client.Email,
		CreatedAt: timestamppb.New(client.CreatedAt),
	}
}
