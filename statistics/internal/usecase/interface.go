package usecase

import (
	"context"
	"taskManager/statistics/internal/model"
)

type StatisticsRepo interface {
	GetUserStatistics(ctx context.Context) (*model.UserStatistics, error)
	GetTaskStatistics(ctx context.Context) (*[]model.TaskStatistics, error)
	GetTaskStatisticsByUserID(ctx context.Context, userID int64) (*model.TaskStatistics, error)
	CreateUserStatistics(ctx context.Context, client model.UserNats) error
	CreateTaskStatistics(ctx context.Context, client model.TaskNats) error
}
