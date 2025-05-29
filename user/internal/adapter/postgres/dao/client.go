package dao

import (
	"taskManager/user/internal/model"
	"time"
)

type Client struct {
	ID        int64     `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"size:50;not null"`
	Password  string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:254;not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt bool      `gorm:"index"`
}

func (c *Client) TableName() string {
	return "users"
}

func FromClient(user *model.User) *Client {
	return &Client{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToClient(user *Client) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
