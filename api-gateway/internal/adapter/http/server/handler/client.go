package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"taskManager/api-gateway/config"
	"taskManager/proto/gen/statistics"
	taskSvc "taskManager/proto/gen/task"
	userSvc "taskManager/proto/gen/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Client struct {
	userClient       userSvc.UserServiceClient
	taskClient       taskSvc.TaskServiceClient
	statisticsClient statistics.StatisticsServiceClient
}

func New(cfg config.Server, userConn *grpc.ClientConn, taskConn *grpc.ClientConn, statisticsConn *grpc.ClientConn) *Client {
	return &Client{
		userClient:       userSvc.NewUserServiceClient(userConn),
		taskClient:       taskSvc.NewTaskServiceClient(taskConn),
		statisticsClient: statistics.NewStatisticsServiceClient(statisticsConn),
	}
}

func (c *Client) Register(ctx *gin.Context) {
	var req userSvc.UserCreateRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	user, err := c.userClient.UserCreate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *Client) Login(ctx *gin.Context) {
	_, err := ctx.Cookie("auth_token")
	if err == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "already logged in"})
		return
	}

	var req userSvc.UserLoginRequest
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.userClient.UserLogin(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie(
		"auth_token",
		resp.Token,
		3600*1,
		"/",
		"",
		false,
		false,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"id":    resp.Id,
		"token": resp.Token,
	})
}

func (c *Client) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	resp, err := c.userClient.UserGet(ctx, &userSvc.UserIDRequest{Id: idInt})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": resp})
}

func (c *Client) DeleteUserByID(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	idInt, _ := strconv.ParseInt(strconv.FormatInt(userID.(int64), 10), 10, 64)

	_, err := c.userClient.UserDelete(ctx, &userSvc.UserIDRequest{Id: idInt})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie(
		"auth_token", // имя куки
		"",           // пустое значение
		-1,           // MaxAge < 0 означает немедленное удаление
		"/",          // путь
		"",           // домен (оставьте пустым для текущего домена)
		false,        // secure
		true,         // httpOnly
	)

	ctx.JSON(http.StatusOK, gin.H{"user": "deleted successfully!"})
}

func (c *Client) CreateTask(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	req := &taskSvc.CreateUserTaskRequest{
		UserId: fmt.Sprintf("%d", userID),
	}

	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.taskClient.CreateUserTask(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Created": resp.Task})
}

func (c *Client) GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	req := &taskSvc.GetTaskByIDRequest{
		TaskId: id,
	}

	task, err := c.taskClient.GetTaskByID(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

func (c *Client) GetUserTaskByID(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	req := &taskSvc.GetUserTaskByIDRequest{
		UserId: fmt.Sprintf("%d", userID),
		TaskId: fmt.Sprintf("%d", idInt),
	}

	resp, err := c.taskClient.GetUserTaskByID(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": resp.Task})
}

func (c *Client) GetAllUserTask(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	resp, err := c.taskClient.GetUserAllTask(ctx, &taskSvc.GetUserAllTaskRequest{UserId: fmt.Sprintf("%d", userID)})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tasks": resp.Task})
}

func (c *Client) DeleteUserTaskByID(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	fmt.Printf("userID type: %T, value: %v\n", userID, userID)

	id := ctx.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	var userIDInt64 int64
	switch v := userID.(type) {
	case int64:
		userIDInt64 = v
	case int:
		userIDInt64 = int64(v)
	case float64:
		userIDInt64 = int64(v)
	case string:
		var parseErr error
		userIDInt64, parseErr = strconv.ParseInt(v, 10, 64)
		if parseErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user ID"})
			return
		}
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected user ID type"})
		return
	}

	req := &taskSvc.DeleteUserTaskByIDRequest{
		UserId: fmt.Sprintf("%d", userIDInt64),
		TaskId: fmt.Sprintf("%d", idInt),
	}

	fmt.Printf("Delete request - UserID: %s, TaskID: %s\n", req.UserId, req.TaskId)

	_, err = c.taskClient.DeleteUserTaskByID(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"deleted": true})
}

func (c *Client) GetUserStatistics(ctx *gin.Context) {
	var req statistics.Empty

	returnThing, err := c.statisticsClient.GetUserStatistics(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"statistics": returnThing})
}

func (c *Client) GetTaskStatistics(ctx *gin.Context) {
	var req statistics.Empty

	returnThing, err := c.statisticsClient.GetTaskStatistics(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"statistics": returnThing})
}

func (c *Client) GetTaskStatisticsByUserID(ctx *gin.Context) {

	id := ctx.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	var req statistics.TaskByIDRequest

	req.Id = idInt

	returnThing, err := c.statisticsClient.GetTaskStatisticsByUserID(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"statistics": returnThing})
}

func (c *Client) UpdateUserTask(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	idInt, _ := strconv.ParseInt(strconv.FormatInt(userID.(int64), 10), 10, 64)

	var req taskSvc.UpdateUserTaskRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.UserId = fmt.Sprintf("%d", idInt)

	updatedTask, err := c.taskClient.UpdateUserTask(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"task": updatedTask})
}
