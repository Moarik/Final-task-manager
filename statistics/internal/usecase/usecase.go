package usecase

import (
	"context"
	"taskManager/statistics/internal/model"
)

type Statistics struct {
	Repo StatisticsRepo
}

func New(repo StatisticsRepo) *Statistics {
	return &Statistics{
		Repo: repo,
	}
}

func (s *Statistics) GetUserStatistics(ctx context.Context) (*model.UserStatistics, error) {
	statistics, err := s.Repo.GetUserStatistics(ctx)
	if err != nil {
		return nil, err
	}

	return statistics, nil
}

func (s *Statistics) GetTaskStatistics(ctx context.Context) (*[]model.TaskStatistics, error) {
	statistics, err := s.Repo.GetTaskStatistics(ctx)
	if err != nil {
		return nil, err
	}

	return statistics, nil
}

func (s *Statistics) GetTaskStatisticsByUserID(ctx context.Context, userID int64) (*model.TaskStatistics, error) {
	statistics, err := s.Repo.GetTaskStatisticsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return statistics, nil
}

func (s *Statistics) CreateUserStatistics(ctx context.Context, client model.UserNats) error {
	err := s.Repo.CreateUserStatistics(ctx, client)
	if err != nil {
		return err
	}

	return nil
}

func (s *Statistics) CreateTaskStatistics(ctx context.Context, client model.TaskNats) error {
	err := s.Repo.CreateTaskStatistics(ctx, client)
	if err != nil {
		return err
	}

	return nil
}
