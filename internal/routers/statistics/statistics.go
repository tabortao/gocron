package statistics

import (
	"github.com/gin-gonic/gin"
	"github.com/tabortao/gocron/internal/models"
	"github.com/tabortao/gocron/internal/modules/logger"
	"github.com/tabortao/gocron/internal/modules/utils"
	"github.com/tabortao/gocron/internal/routers/base"
)

// OverviewData 概览统计数据
type OverviewData struct {
	TotalTasks      int64               `json:"total_tasks"`
	TodayExecutions int64               `json:"today_executions"`
	SuccessRate     float64             `json:"success_rate"`
	FailedCount     int64               `json:"failed_count"`
	Last7Days       []models.DailyStats `json:"last_7_days"`
}

// Overview 获取统计概览数据
func Overview(c *gin.Context) {
	taskModel := models.Task{}
	taskLogModel := models.TaskLog{}

	// 1. 获取启用的任务总数
	totalTasks, err := taskModel.Total(models.CommonMap{"Status": int(models.Enabled)})
	if err != nil {
		logger.Error("Failed to get total tasks:", err)
		base.RespondError(c, "Failed to get total tasks", err)
		return
	}

	// 2. 获取今日统计数据
	todayTotal, todaySuccess, todayFailed, err := taskLogModel.GetTodayStats()
	if err != nil {
		logger.Error("Failed to get today's statistics:", err)
		base.RespondError(c, "Failed to get today's statistics", err)
		return
	}

	// 3. 计算成功率
	var successRate float64
	if todayTotal > 0 {
		successRate = float64(todaySuccess) / float64(todayTotal) * 100
		// 保留1位小数
		successRate = float64(int(successRate*10)) / 10
	}

	// 4. 获取最近7天趋势
	last7Days, err := taskLogModel.GetLast7DaysTrend()
	if err != nil {
		logger.Error("Failed to get trend data:", err)
		base.RespondError(c, "Failed to get trend data", err)
		return
	}

	// 组装返回数据
	data := OverviewData{
		TotalTasks:      totalTasks,
		TodayExecutions: todayTotal,
		SuccessRate:     successRate,
		FailedCount:     todayFailed,
		Last7Days:       last7Days,
	}

	base.RespondSuccess(c, utils.SuccessContent, data)
}
