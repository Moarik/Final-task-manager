package handler

import (
	"context"
	"taskManager/statistics/internal/model"
)

type StatisticsUsecase interface {
	CreateUserStatistics(ctx context.Context, client model.UserNats) error
	CreateTaskStatistics(ctx context.Context, client model.TaskNats) error
}
