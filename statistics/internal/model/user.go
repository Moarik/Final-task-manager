package model

type UserStatistics struct {
	ID         int64 `json:"id"`
	TotalUsers int   `json:"total_users"`
}

type UserNats struct {
	ID int64
}
