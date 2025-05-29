package dto

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"taskManager/user/internal/model"
)

type ClientCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ClientLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ClientCreateResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func ToClientFromLoginRequest(ctx *gin.Context) (*model.User, error) {
	var req ClientLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, err
	}

	return &model.User{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func ToClientFromCreateRequest(ctx *gin.Context) (model.User, error) {
	var req ClientCreateRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return model.User{}, err
	}

	if req.Name == "" {
		return model.User{}, errors.New("name is required")
	}

	err = validateClientCreateRequest(req)
	if err != nil {
		return model.User{}, err
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: passwordHash,
	}, nil
}

func FromClientToCreateResponse(client model.User) ClientCreateResponse {
	return ClientCreateResponse{
		ID:    client.ID,
		Email: client.Email,
	}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
