package task

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	taskSvc "taskManager/proto/gen/task"
	"taskManager/task/internal/adapter/grpc/service/task/dto"
)

type Client struct {
	taskSvc.UnimplementedTaskServiceServer

	uc TaskUsecase
}

func New(uc TaskUsecase) *Client {
	return &Client{
		uc: uc,
	}
}

func (c *Client) CreateUserTask(ctx context.Context, req *taskSvc.CreateUserTaskRequest) (*taskSvc.TaskResponse, error) {
	task, err := dto.ToTaskFromRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	returnTask, err := c.uc.CreateUserTaskService(ctx, task)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sendTask := &taskSvc.Task{
		TaskId:      strconv.FormatInt(returnTask.ID, 10),
		UserId:      strconv.FormatInt(returnTask.UserID, 10),
		Title:       returnTask.Title,
		Description: returnTask.Description,
		IsCompleted: returnTask.Completed,
	}

	return &taskSvc.TaskResponse{
		Task: sendTask,
	}, nil
}

func (c *Client) GetUserTaskByID(ctx context.Context, req *taskSvc.GetUserTaskByIDRequest) (*taskSvc.TaskResponse, error) {
	userIdConverted, _ := strconv.ParseInt(req.UserId, 10, 64)
	taskIdConverted, _ := strconv.ParseInt(req.TaskId, 10, 64)

	task, err := c.uc.GetUserTaskByIDService(ctx, userIdConverted, taskIdConverted)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sendTask := &taskSvc.Task{
		TaskId:      strconv.FormatInt(task.ID, 10),
		UserId:      strconv.FormatInt(task.UserID, 10),
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.Completed,
	}

	return &taskSvc.TaskResponse{
		Task: sendTask,
	}, nil
}

func (c *Client) GetUserAllTask(ctx context.Context, req *taskSvc.GetUserAllTaskRequest) (*taskSvc.TasksResponse, error) {
	userIdConverted, _ := strconv.ParseInt(req.UserId, 10, 64)

	tasks, err := c.uc.GetAllUserTasksByIDService(ctx, userIdConverted)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var sendTasks []*taskSvc.Task
	for _, task := range tasks {
		sendTask := &taskSvc.Task{
			TaskId:      strconv.FormatInt(task.ID, 10),
			UserId:      strconv.FormatInt(task.UserID, 10),
			Title:       task.Title,
			Description: task.Description,
			IsCompleted: task.Completed,
		}

		sendTasks = append(sendTasks, sendTask)
	}

	return &taskSvc.TasksResponse{
		Task: sendTasks,
	}, nil
}

func (c *Client) GetTaskByID(ctx context.Context, req *taskSvc.GetTaskByIDRequest) (*taskSvc.TaskResponse, error) {
	taskIDConverted, _ := strconv.ParseInt(req.TaskId, 10, 64)

	task, err := c.uc.GetTaskByIDService(ctx, taskIDConverted)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sendTask := &taskSvc.Task{
		TaskId:      strconv.FormatInt(task.ID, 10),
		UserId:      strconv.FormatInt(task.UserID, 10),
		Title:       task.Title,
		Description: task.Description,
		IsCompleted: task.Completed,
	}

	return &taskSvc.TaskResponse{
		Task: sendTask,
	}, nil
}

func (c *Client) DeleteUserTaskByID(ctx context.Context, req *taskSvc.DeleteUserTaskByIDRequest) (*taskSvc.Empty, error) {
	taskIdConverted, _ := strconv.ParseInt(req.TaskId, 10, 64)
	userIdConverted, _ := strconv.ParseInt(req.UserId, 10, 64)

	err := c.uc.DeleteUserTaskByIDService(ctx, taskIdConverted, userIdConverted)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &taskSvc.Empty{}, nil
}

func (c *Client) UpdateUserTask(ctx context.Context, req *taskSvc.UpdateUserTaskRequest) (*taskSvc.TaskResponse, error) {
	task, err := dto.ToTaskFromRequestUpdate(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updatedTask, err := c.uc.UpdateUserTaskByIDService(ctx, task)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sendTask := &taskSvc.Task{
		TaskId:      strconv.FormatInt(updatedTask.ID, 10),
		UserId:      strconv.FormatInt(updatedTask.UserID, 10),
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		IsCompleted: updatedTask.Completed,
	}

	return &taskSvc.TaskResponse{
		Task: sendTask,
	}, nil
}
