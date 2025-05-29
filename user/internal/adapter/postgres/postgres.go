package postgres

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"taskManager/user/internal/adapter/postgres/dao"
	"taskManager/user/internal/model"
	"taskManager/user/pkg/crypto"
	"taskManager/user/pkg/postgre"
)

type Client struct {
	DB postgre.Database
}

func NewClient(db postgre.Database) *Client {
	return &Client{DB: db}
}

func (c *Client) Create(ctx context.Context, user *model.User) (model.User, error) {
	daoUser := dao.FromClient(user)

	if err := c.DB.GetDB().WithContext(ctx).Create(&daoUser).Error; err != nil {
		return model.User{}, err
	}

	returnUser := dao.ToClient(daoUser)

	return *returnUser, nil
}

func (c *Client) Login(ctx context.Context, user *model.User) (int64, error) {
	var userModel model.User

	result := c.DB.GetDB().WithContext(ctx).Where("email = ? and deleted_at = false", user.Email).First(&userModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("user with such email %s not found", user.Email)
		}
		return 0, fmt.Errorf("%s", result.Error)
	}

	if !crypto.VerifyPassword(userModel.Password, user.Password) {
		return 0, fmt.Errorf("wrong password")
	}

	return userModel.ID, nil
}

func (c *Client) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var daoUser dao.Client
	err := c.DB.GetDB().WithContext(ctx).Where("deleted_at = false").First(&daoUser, id).Error
	if err != nil {
		return nil, err
	}

	return dao.ToClient(&daoUser), nil
}

func (c *Client) DeleteByID(ctx context.Context, id int64) error {
	return c.DB.GetDB().WithContext(ctx).
		Model(&dao.Client{}).
		Where("id = ?", id).
		Update("deleted_at", true).Error
}
