package tasklog

// 任务日志

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tabortao/gocron/internal/models"
	"github.com/tabortao/gocron/internal/modules/i18n"
	"github.com/tabortao/gocron/internal/modules/logger"
	rpcClient "github.com/tabortao/gocron/internal/modules/rpc/client"
	"github.com/tabortao/gocron/internal/modules/utils"
	"github.com/tabortao/gocron/internal/routers/base"
	"github.com/tabortao/gocron/internal/service"
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
	base.RespondSuccess(c, utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  logs,
	})
}

// 清空日志
func Clear(c *gin.Context) {
	taskLogModel := new(models.TaskLog)
	_, err := taskLogModel.Clear()
	if err != nil {
		base.RespondErrorWithDefaultMsg(c, err)
	} else {
		base.RespondSuccessWithDefaultMsg(c, nil)
	}
}

func Output(c *gin.Context) {
	logId, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || logId <= 0 {
		base.RespondError(c, i18n.T(c, "invalid_log_id"))
		return
	}

	taskLogModel := new(models.TaskLog)
	taskLog, err := taskLogModel.Detail(logId)
	if err != nil {
		base.RespondErrorWithDefaultMsg(c, err)
		return
	}
	if taskLog.Id <= 0 {
		base.RespondError(c, i18n.T(c, "invalid_log_id"))
		return
	}

	if taskLog.Protocol != models.TaskRPC || taskLog.Status != models.Running {
		base.RespondSuccess(c, utils.SuccessContent, map[string]interface{}{
			"output": taskLog.Result,
			"status": taskLog.Status,
		})
		return
	}

	taskModel := new(models.Task)
	task, err := taskModel.Detail(taskLog.TaskId)
	if err != nil {
		base.RespondError(c, i18n.T(c, "get_task_info_failed")+"#"+err.Error(), err)
		return
	}
	if len(task.Hosts) == 0 {
		base.RespondError(c, i18n.T(c, "task_node_list_empty"))
		return
	}

	aggregationResult := ""
	for _, host := range task.Hosts {
		output, tailErr := rpcClient.Tail(host.Name, host.Port, logId)
		errorMessage := ""
		if tailErr != nil {
			errorMessage = strings.TrimSpace(tailErr.Error()) + "\n"
		}
		output = strings.TrimSpace(output)
		outputMessage := fmt.Sprintf("Host: [%s-%s:%d]\n%s%s\n", host.Alias, host.Name, host.Port, errorMessage, output)
		aggregationResult += outputMessage
	}

	base.RespondSuccess(c, utils.SuccessContent, map[string]interface{}{
		"output": aggregationResult,
		"status": taskLog.Status,
	})
}

// 停止运行中的任务
func Stop(c *gin.Context) {
	id, err := strconv.ParseInt(c.PostForm("id"), 10, 64)
	if err != nil || id <= 0 {
		base.RespondError(c, i18n.T(c, "invalid_log_id"))
		return
	}
	taskId, err := strconv.Atoi(c.PostForm("task_id"))
	if err != nil || taskId <= 0 {
		base.RespondError(c, i18n.T(c, "invalid_task_id"))
		return
	}
	taskModel := new(models.Task)
	task, err := taskModel.Detail(taskId)
	if err != nil {
		base.RespondError(c, i18n.T(c, "get_task_info_failed")+"#"+err.Error(), err)
		return
	}
	if task.Protocol != models.TaskRPC {
		base.RespondError(c, i18n.T(c, "only_shell_task_can_stop"))
		return
	}
	if len(task.Hosts) == 0 {
		base.RespondError(c, i18n.T(c, "task_node_list_empty"))
		return
	}
	for _, host := range task.Hosts {
		service.ServiceTask.Stop(host.Name, host.Port, id)
	}

	base.RespondSuccess(c, i18n.T(c, "stop_task_sent"), nil)
}

// 删除N个月前的日志
func Remove(c *gin.Context) {
	month, _ := strconv.Atoi(c.Param("id"))
	if month < 1 || month > 12 {
		base.RespondError(c, i18n.T(c, "param_range_1_12"))
		return
	}
	taskLogModel := new(models.TaskLog)
	_, err := taskLogModel.Remove(month)
	if err != nil {
		base.RespondError(c, i18n.T(c, "delete_failed"), err)
	} else {
		base.RespondSuccess(c, i18n.T(c, "delete_success"), nil)
	}
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
