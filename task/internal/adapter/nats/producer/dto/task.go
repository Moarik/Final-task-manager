package dto

import (
	"taskManager/proto/gen/statistics"
	"taskManager/task/internal/model"
)

func FromTask(task model.Task) statistics.TaskCreateNats {
	return statistics.TaskCreateNats{
		Id: task.UserID,
	}
}
