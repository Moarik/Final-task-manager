package postgres

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"taskManager/statistics/internal/adapter/postgres/dao"
	"taskManager/statistics/internal/model"
	"taskManager/statistics/pkg/postgre"
)

type Statistics struct {
	DB postgre.Database
}

func New(db postgre.Database) *Statistics {
	return &Statistics{
		DB: db,
	}
}

func (s *Statistics) GetUserStatistics(ctx context.Context) (*model.UserStatistics, error) {
	var daoUser dao.UserStatistics

	err := s.DB.GetDB().WithContext(ctx).Unscoped().First(&daoUser).Error
	if err != nil {
		return nil, err
	}

	returnModel := &model.UserStatistics{
		ID:         daoUser.ID,
		TotalUsers: daoUser.TotalUsers,
	}

	return returnModel, nil
}

func (s *Statistics) GetTaskStatisticsByUserID(ctx context.Context, userID int64) (*model.TaskStatistics, error) {
	var daoTask dao.TaskStatistics

	err := s.DB.GetDB().WithContext(ctx).
		Unscoped().
		Where("id = ?", userID).
		First(&daoTask).Error
	if err != nil {
		return nil, err
	}

	returnModel := &model.TaskStatistics{
		ID:         daoTask.ID,
		TotalTasks: daoTask.TotalTasks,
	}

	return returnModel, nil
}

func (s *Statistics) GetTaskStatistics(ctx context.Context) (*[]model.TaskStatistics, error) {
	var daoTasks []dao.TaskStatistics

	err := s.DB.GetDB().
		WithContext(ctx).
		Table("tasks").
		Select("id, COUNT(*) as task_count").
		Group("id").
		Scan(&daoTasks).Error
	if err != nil {
		return nil, err
	}

	var stats []model.TaskStatistics
	for _, r := range daoTasks {
		stats = append(stats, model.TaskStatistics{
			ID:         r.ID,
			TotalTasks: r.TotalTasks,
		})
	}

	return &stats, nil
}

func (s *Statistics) CreateUserStatistics(ctx context.Context, client model.UserNats) error {
	if client.ID == 0 {
		return errors.New("invalid user ID: 0")
	}

	var stat dao.UserStatistics
	db := s.DB.GetDB().WithContext(ctx)

	err := db.First(&stat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		stat = dao.UserStatistics{
			TotalUsers: 1,
		}
		return db.Create(&stat).Error
	} else if err != nil {
		return err
	}

	stat.TotalUsers += 1
	return db.Save(&stat).Error
}

func (s *Statistics) CreateTaskStatistics(ctx context.Context, client model.TaskNats) error {
	if client.ID == 0 {
		return errors.New("invalid task ID: 0")
	}

	fmt.Println("TASK FROM STATISTICS REPOSITORY:", client)

	var stat dao.TaskStatistics
	db := s.DB.GetDB().WithContext(ctx)

	err := db.First(&stat, client.ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		stat = dao.TaskStatistics{
			ID:         client.ID,
			TotalTasks: 1,
		}
		return db.Create(&stat).Error
	} else if err != nil {
		return err
	}

	stat.TotalTasks += 1
	return db.Save(&stat).Error
}
