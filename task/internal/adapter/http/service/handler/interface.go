package handler

import (
	"context"
	"taskManager/task/internal/model"
)

type TaskUsecase interface {
	CreateUserTaskService(ctx context.Context, task model.Task) (*model.Task, error)
	GetUserTaskByIDService(ctx context.Context, userId, taskId int64) (*model.Task, error)
	DeleteUserTaskByIDService(ctx context.Context, userId, taskId int64) error
	UpdateUserTaskByIDService(ctx context.Context, task model.Task) (*model.Task, error)
	GetAllUserTasksByIDService(ctx context.Context, id int64) ([]*model.Task, error)
	GetTaskByIDService(ctx context.Context, id int64) (*model.Task, error)
}
