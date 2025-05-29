package statistics

import (
	"context"
	"taskManager/statistics/internal/model"
)

type StatisticsUsecase interface {
	GetUserStatistics(ctx context.Context) (*model.UserStatistics, error)
	GetTaskStatistics(ctx context.Context) (*[]model.TaskStatistics, error)
	GetTaskStatisticsByUserID(ctx context.Context, userID int64) (*model.TaskStatistics, error)
}
