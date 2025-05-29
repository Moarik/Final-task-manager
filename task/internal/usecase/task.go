package usecase

import (
	"context"
	"fmt"
	"taskManager/task/internal/model"
)

type Task struct {
	Repo     TaskRepo
	Producer TaskEventStorage
	Redis    RedisCache
}

func NewTask(repo TaskRepo, producer TaskEventStorage, redis RedisCache) *Task {
	return &Task{
		Repo:     repo,
		Producer: producer,
		Redis:    redis,
	}
}

func (t *Task) CreateUserTaskService(ctx context.Context, task model.Task) (*model.Task, error) {
	taskModel, err := t.Repo.CreateUserTask(ctx, task)
	if err != nil {
		return nil, err
	}

	err = t.Producer.Push(ctx, *taskModel)
	if err != nil {
		return nil, err
	}

	err = t.Redis.Set(ctx, *taskModel)
	if err != nil {
		return nil, err
	}

	return taskModel, nil
}

func (t *Task) GetUserTaskByIDService(ctx context.Context, userId, taskId int64) (*model.Task, error) {
	taskModel, err := t.Repo.GetUserTaskByID(ctx, userId, taskId)
	if err != nil {
		return nil, err
	}

	return taskModel, nil
}

func (t *Task) DeleteUserTaskByIDService(ctx context.Context, userID, taskID int64) error {
	err := t.Repo.DeleteUserTaskByID(ctx, userID, taskID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) UpdateUserTaskByIDService(ctx context.Context, task model.Task) (*model.Task, error) {
	updatedUserTask, err := t.Repo.UpdateUserTaskByID(ctx, task)
	if err != nil {
		return nil, err
	}

	return updatedUserTask, nil
}

func (t *Task) GetAllUserTasksByIDService(ctx context.Context, id int64) ([]*model.Task, error) {
	userTasks, err := t.Repo.GetAllUsersTasks(ctx, id)
	if err != nil {
		return nil, err
	}

	return userTasks, nil
}

func (t *Task) GetTaskByIDService(ctx context.Context, id int64) (*model.Task, error) {
	taskRedis, err := t.Redis.Get(ctx, id)
	if err == nil {
		fmt.Println("GOT FROM CACHE REDIS")
		return &taskRedis, nil
	}

	task, err := t.Repo.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	fmt.Println("GOT FROM DATABASE")
	return task, nil
}

func (t *Task) RefreshCache(ctx context.Context) error {
	tasks, err := t.Repo.GetAllTasks(ctx)
	if err != nil {
		return err
	}

	err = t.Redis.SetMany(ctx, tasks)
	if err != nil {
		return nil
	}

	return nil
}
