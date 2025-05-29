package usecase

import (
	"context"
	"taskManager/user/internal/model"
)

type ClientRepo interface {
	Create(ctx context.Context, user *model.User) (model.User, error)
	Login(ctx context.Context, user *model.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	DeleteByID(ctx context.Context, id int64) error
}

type ClientEventStorage interface {
	Push(ctx context.Context, user model.User) error
}
