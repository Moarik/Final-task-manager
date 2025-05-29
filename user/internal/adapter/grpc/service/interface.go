package service

import (
	"context"
	"taskManager/user/internal/model"
)

type UserUsecase interface {
	Create(ctx context.Context, request model.User) (model.User, error)
	Login(ctx context.Context, request model.User) (int64, error)
	Get(ctx context.Context, id int64) (model.User, error)
	Delete(ctx context.Context, id int64) error
}
