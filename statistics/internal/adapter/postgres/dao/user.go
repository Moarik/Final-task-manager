package dao

type UserStatistics struct {
	ID         int64 `gorm:"column:id"`
	TotalUsers int   `gorm:"column:total_users"`
}
