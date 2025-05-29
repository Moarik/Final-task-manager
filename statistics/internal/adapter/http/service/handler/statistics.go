package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Statistics struct {
	uc StatisticsUsecase
}

func New(uc StatisticsUsecase) *Statistics {
	return &Statistics{
		uc: uc,
	}
}

func (s *Statistics) GetUserStatistics(ctx *gin.Context) {
	userStats, err := s.uc.GetUserStatistics(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userStats)
}

func (s *Statistics) GetTaskStatistics(ctx *gin.Context) {
	taskStats, err := s.uc.GetTaskStatistics(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, taskStats)
}
