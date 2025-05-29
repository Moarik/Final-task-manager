package dto

import (
	userSvc "taskManager/proto/gen/user"
	"taskManager/user/internal/model"
)

func ToClientFromRequest(req *userSvc.UserCreateRequest) (model.User, error) {
	return model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func ToClientFromLogin(req *userSvc.UserLoginRequest) (model.User, error) {
	return model.User{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}
