package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskManager/user/internal/adapter/http/service/handler/dto"
	"taskManager/user/pkg/security"
)

type Client struct {
	uc UserUsecase
}

func NewClient(uc UserUsecase) *Client {
	return &Client{
		uc: uc,
	}
}

func (c *Client) Create(ctx *gin.Context) {
	client, err := dto.ToClientFromCreateRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUser, err := c.uc.Create(ctx, client)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, dto.FromClientToCreateResponse(newUser))
}

func (c *Client) Login(ctx *gin.Context) {
	client, err := dto.ToClientFromLoginRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := c.uc.Login(ctx, *client)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := security.CreateToken(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *Client) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	client, err := c.uc.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, client)
}

func (c *Client) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = c.uc.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
