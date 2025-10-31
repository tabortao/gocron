package tasklog

// 任务日志

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/models"
	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/gocronx-team/gocron/internal/modules/utils"
	"github.com/gocronx-team/gocron/internal/routers/base"
	"github.com/gocronx-team/gocron/internal/service"
)

func Index(c *gin.Context) {
	logModel := new(models.TaskLog)
	queryParams := parseQueryParams(c)
	total, err := logModel.Total(queryParams)
	if err != nil {
		logger.Error(err)
	}
	logs, err := logModel.List(queryParams)
	if err != nil {
		logger.Error(err)
	}
	jsonResp := utils.JsonResponse{}
	result := jsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  logs,
	})
	c.String(http.StatusOK, result)
}

// 清空日志
func Clear(c *gin.Context) {
	taskLogModel := new(models.TaskLog)
	_, err := taskLogModel.Clear()
	json := utils.JsonResponse{}
	var result string
	if err != nil {
		result = json.CommonFailure(utils.FailureContent)
	} else {
		result = json.Success(utils.SuccessContent, nil)
	}
	c.String(http.StatusOK, result)
}

// 停止运行中的任务
func Stop(c *gin.Context) {
	var form struct {
		Id     int64 `form:"id" binding:"required"`
		TaskId int   `form:"task_id" binding:"required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		json := utils.JsonResponse{}
		result := json.CommonFailure("参数错误")
		c.String(http.StatusOK, result)
		return
	}
	id := form.Id
	taskId := form.TaskId
	taskModel := new(models.Task)
	task, err := taskModel.Detail(taskId)
	json := utils.JsonResponse{}
	var result string
	if err != nil {
		result = json.CommonFailure("获取任务信息失败#"+err.Error(), err)
		c.String(http.StatusOK, result)
		return
	}
	if task.Protocol != models.TaskRPC {
		result = json.CommonFailure("仅支持SHELL任务手动停止")
		c.String(http.StatusOK, result)
		return
	}
	if len(task.Hosts) == 0 {
		result = json.CommonFailure("任务节点列表为空")
		c.String(http.StatusOK, result)
		return
	}
	for _, host := range task.Hosts {
		service.ServiceTask.Stop(host.Name, host.Port, id)
	}

	result = json.Success("已执行停止操作, 请等待任务退出", nil)
	c.String(http.StatusOK, result)
}

// 删除N个月前的日志
func Remove(c *gin.Context) {
	month, _ := strconv.Atoi(c.Param("id"))
	json := utils.JsonResponse{}
	var result string
	if month < 1 || month > 12 {
		result = json.CommonFailure("参数取值范围1-12")
		c.String(http.StatusOK, result)
		return
	}
	taskLogModel := new(models.TaskLog)
	_, err := taskLogModel.Remove(month)
	if err != nil {
		result = json.CommonFailure("删除失败", err)
	} else {
		result = json.Success("删除成功", nil)
	}
	c.String(http.StatusOK, result)
}

// 解析查询参数
func parseQueryParams(c *gin.Context) models.CommonMap {
	var params models.CommonMap = models.CommonMap{}
	taskId, _ := strconv.Atoi(c.Query("task_id"))
	protocol, _ := strconv.Atoi(c.Query("protocol"))
	status, _ := strconv.Atoi(c.Query("status"))
	params["TaskId"] = taskId
	params["Protocol"] = protocol
	if status >= 0 {
		status -= 1
	}
	params["Status"] = status
	base.ParsePageAndPageSize(c, params)

	return params
}
