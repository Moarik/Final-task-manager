package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskManager/task/internal/adapter/http/service/handler/dto"
)

type Task struct {
	uc TaskUsecase
}

func New(uc TaskUsecase) *Task {
	return &Task{uc: uc}
}

func (t *Task) CreateUserTask(ctx *gin.Context) {
	reqModel, err := dto.FromCreateRequestToTask(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskModel, err := t.uc.CreateUserTaskService(ctx.Request.Context(), *reqModel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, taskModel)
}

func (t *Task) GetUserTaskByID(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	userId := 1

	task, err := t.uc.GetUserTaskByIDService(ctx.Request.Context(), int64(userId), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (t *Task) DeleteUserTaskByID(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	userID, _ := strconv.ParseInt("1", 10, 64)

	err := t.uc.DeleteUserTaskByIDService(ctx, id, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	return
}

func (t *Task) UpdateUserTaskByID(ctx *gin.Context) {
	reqModel, err := dto.FromUpdateRequestToTask(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUserTask, err := t.uc.UpdateUserTaskByIDService(ctx.Request.Context(), *reqModel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedUserTask)
}

func (t *Task) GetAllUserTasksByID(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	tasks, err := t.uc.GetAllUserTasksByIDService(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

func (t *Task) GetTaskByID(ctx *gin.Context) {
	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	task, err := t.uc.GetTaskByIDService(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, task)
}
