package model

import "time"

type Task struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskRedis struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
