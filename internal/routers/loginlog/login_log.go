package loginlog

import (
	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/models"
	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/gocronx-team/gocron/internal/modules/utils"
	"github.com/gocronx-team/gocron/internal/routers/base"
)

func Index(c *gin.Context) {
	loginLogModel := new(models.LoginLog)
	params := models.CommonMap{}
	base.ParsePageAndPageSize(c, params)
	total, err := loginLogModel.Total()
	if err != nil {
		logger.Error(err)
	}
	loginLogs, err := loginLogModel.List(params)
	if err != nil {
		logger.Error(err)
	}

	base.RespondSuccess(c, utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  loginLogs,
	})
}
