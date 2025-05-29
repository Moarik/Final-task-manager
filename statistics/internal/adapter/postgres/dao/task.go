package dao

type TaskStatistics struct {
	ID         int64 `gorm:"column:id"`
	TotalTasks int   `gorm:"column:total_tasks"`
}
