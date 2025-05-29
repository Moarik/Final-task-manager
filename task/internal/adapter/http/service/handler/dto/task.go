package dto

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"taskManager/task/internal/model"
)

type TaskCreateRequest struct {
	UserID      int64  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskUpdateRequest struct {
	UserID      *int64  `json:"user_id,omitempty"`
	TaskID      *int64  `json:"task_id,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	IsCompleted *bool   `json:"is_completed,omitempty"`
}

func (t *TaskCreateRequest) Validate() error {
	if t.Title == "" {
		return errors.New("title is required")
	}

	if t.Description == "" {
		return errors.New("description is required")
	}

	return nil
}

func FromCreateRequestToTask(ctx *gin.Context) (*model.Task, error) {
	var req TaskCreateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	return &model.Task{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
	}, nil
}

func FromUpdateRequestToTask(ctx *gin.Context) (*model.Task, error) {
	var req TaskUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	if req.TaskID == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "task_id is required"})
		return nil, errors.New("task_id is required")
	}

	task := &model.Task{
		ID: *req.TaskID,
	}

	if req.UserID != nil {
		task.UserID = *req.UserID
	}

	if req.IsCompleted != nil {
		task.Completed = *req.IsCompleted
	}

	if req.Description != nil {
		task.Description = *req.Description
	}

	if req.Title != nil {
		task.Title = *req.Title
	}

	return task, nil
}
