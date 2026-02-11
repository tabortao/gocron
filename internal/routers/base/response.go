package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tabortao/gocron/internal/modules/logger"
	"github.com/tabortao/gocron/internal/modules/utils"
)

// RespondSuccess 返回成功响应
func RespondSuccess(c *gin.Context, message string, data interface{}) {
	json := utils.JsonResponse{}
	result := json.Success(message, data)
	c.String(http.StatusOK, result)
}

// RespondSuccessWithDefaultMsg 返回成功响应（使用默认消息）
func RespondSuccessWithDefaultMsg(c *gin.Context, data interface{}) {
	json := utils.JsonResponse{}
	result := json.Success(utils.SuccessContent, data)
	c.String(http.StatusOK, result)
}

// RespondError 返回错误响应
func RespondError(c *gin.Context, message string, err ...error) {
	json := utils.JsonResponse{}
	if len(err) > 0 && err[0] != nil {
		logger.Error(err[0])
	}
	result := json.CommonFailure(message)
	c.String(http.StatusOK, result)
}

// RespondErrorWithDefaultMsg 返回错误响应（使用默认消息）
func RespondErrorWithDefaultMsg(c *gin.Context, err ...error) {
	json := utils.JsonResponse{}
	if len(err) > 0 && err[0] != nil {
		logger.Error(err[0])
	}
	result := json.CommonFailure(utils.FailureContent)
	c.String(http.StatusOK, result)
}

// RespondValidationError 返回表单验证错误响应
func RespondValidationError(c *gin.Context, err error) {
	json := utils.JsonResponse{}
	result := json.CommonFailure(utils.FailureContent, err)
	c.String(http.StatusOK, result)
}

// RespondAuthError 返回认证错误响应
func RespondAuthError(c *gin.Context, message string) {
	json := utils.JsonResponse{}
	result := json.Failure(utils.AuthError, message)
	c.String(http.StatusOK, result)
}
