package dao

import (
	"taskManager/task/internal/model"
	"time"
)

type Task struct {
	ID          int64     `gorm:"type:uuid;primary_key"`
	UserID      int64     `gorm:"index"`
	Title       string    `gorm:"size:50;not null"`
	Description string    `gorm:"size:255;not null"`
	Completed   bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (t *Task) TableName() string {
	return "tasks"
}

func FromTask(task model.Task) *Task {
	return &Task{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func FromDao(task Task) *model.Task {
	return &model.Task{
		ID:          task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func FromDaos(task []Task) []*model.Task {
	result := make([]*model.Task, len(task))
	for i, t := range task {
		result[i] = FromDao(t)
	}

	return result
}
