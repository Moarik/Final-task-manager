package dto

import (
	"strconv"
	taskSvc "taskManager/proto/gen/task"
	"taskManager/task/internal/model"
)

func ToTaskFromRequest(req *taskSvc.CreateUserTaskRequest) (model.Task, error) {
	id, _ := strconv.ParseInt(req.UserId, 10, 64)
	return model.Task{
		UserID:      id,
		Title:       req.Title,
		Description: req.Description,
	}, nil
}

func ToTaskFromRequestUpdate(req *taskSvc.UpdateUserTaskRequest) (model.Task, error) {
	id, _ := strconv.ParseInt(req.UserId, 10, 64)
	taskId, _ := strconv.ParseInt(req.TaskId, 10, 64)

	return model.Task{
		ID:          taskId,
		UserID:      id,
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.IsCompleted,
	}, nil
}
