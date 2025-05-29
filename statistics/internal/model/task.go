package model

type TaskStatistics struct {
	ID         int64 `json:"id"`
	TotalTasks int   `json:"total_tasks"`
}

type TaskNats struct {
	ID int64
}
