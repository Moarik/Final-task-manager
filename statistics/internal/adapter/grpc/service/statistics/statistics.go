package statistics

import (
	"context"
	"taskManager/proto/gen/statistics"
)

type Client struct {
	statistics.UnimplementedStatisticsServiceServer

	uc StatisticsUsecase
}

func New(uc StatisticsUsecase) *Client {
	return &Client{
		uc: uc,
	}
}

func (c *Client) GetUserStatistics(ctx context.Context, req *statistics.Empty) (*statistics.GetUserStatisticsResponse, error) {
	returnTask, err := c.uc.GetUserStatistics(ctx)
	if err != nil {
		return nil, err
	}

	sendTask := &statistics.GetUserStatisticsResponse{
		Id:         returnTask.ID,
		TotalUsers: int32(returnTask.TotalUsers),
	}

	return sendTask, err
}

func (c *Client) GetAllUserTaskStatistics(ctx context.Context, _ *statistics.Empty) (*statistics.GetTaskStatisticsResponseSpecial, error) {
	stats, err := c.uc.GetTaskStatistics(ctx)
	if err != nil {
		return nil, err
	}

	var response statistics.GetTaskStatisticsResponseSpecial
	for _, stat := range *stats {
		response.Statistics = append(response.Statistics, &statistics.UserTaskCount{
			UserId:    stat.ID,
			TaskCount: int32(stat.TotalTasks),
		})
	}

	return &response, nil
}

func (c *Client) GetTaskStatisticsByUserID(ctx context.Context, req *statistics.TaskByIDRequest) (*statistics.GetTaskStatisticsResponse, error) {
	returnTask, err := c.uc.GetTaskStatisticsByUserID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	sendTask := &statistics.GetTaskStatisticsResponse{
		Id:         returnTask.ID,
		TotalTasks: int32(returnTask.TotalTasks),
	}

	return sendTask, err
}
