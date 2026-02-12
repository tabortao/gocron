package task

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/cron"
	"github.com/tabortao/gocron/internal/models"
	"github.com/tabortao/gocron/internal/modules/i18n"
	"github.com/tabortao/gocron/internal/modules/logger"
	"github.com/tabortao/gocron/internal/modules/utils"
	"github.com/tabortao/gocron/internal/routers/base"
	"github.com/tabortao/gocron/internal/service"
)

type TaskForm struct {
	Id               int                         `form:"id" json:"id"`
	Level            models.TaskLevel            `form:"level" json:"level" binding:"required,oneof=1 2"`
	DependencyStatus models.TaskDependencyStatus `form:"dependency_status" json:"dependency_status" binding:"oneof=1 2"`
	DependencyTaskId string                      `form:"dependency_task_id" json:"dependency_task_id"`
	Name             string                      `form:"name" json:"name" binding:"required,max=32"`
	Spec             string                      `form:"spec" json:"spec"`
	Protocol         models.TaskProtocol         `form:"protocol" json:"protocol" binding:"oneof=1 2"`
	Command          string                      `form:"command" json:"command" binding:"required,max=65535"`
	HttpMethod       models.TaskHTTPMethod       `form:"http_method" json:"http_method" binding:"oneof=1 2"`
	Timeout          int                         `form:"timeout" json:"timeout" binding:"min=0,max=86400"`
	Multi            int8                        `form:"multi" json:"multi" binding:"oneof=0 1"`
	RetryTimes       int8                        `form:"retry_times" json:"retry_times"`
	RetryInterval    int16                       `form:"retry_interval" json:"retry_interval"`
	HostId           string                      `form:"host_id" json:"host_id"`
	Tag              string                      `form:"tag" json:"tag"`
	Remark           string                      `form:"remark" json:"remark"`
	NotifyStatus     int8                        `form:"notify_status" json:"notify_status" binding:"oneof=0 1 2 3"`
	NotifyType       int8                        `form:"notify_type" json:"notify_type"`
	NotifyReceiverId string                      `form:"notify_receiver_id" json:"notify_receiver_id"`
	NotifyKeyword    string                      `form:"notify_keyword" json:"notify_keyword"`
}

// 首页
func Index(c *gin.Context) {
	taskModel := new(models.Task)
	queryParams := parseQueryParams(c)
	total, err := taskModel.Total(queryParams)
	if err != nil {
		logger.Error(err)
	}
	tasks, err := taskModel.List(queryParams)
	if err != nil {
		logger.Error(err)
	}
	for i, item := range tasks {
		tasks[i].NextRunTime = models.NextRunTime(service.ServiceTask.NextRunTime(item))
	}
	jsonResp := utils.JsonResponse{}
	result := jsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  tasks,
	})
	c.String(http.StatusOK, result)
}

// Detail 任务详情
func Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	taskModel := new(models.Task)
	task, err := taskModel.Detail(id)
	jsonResp := utils.JsonResponse{}
	var result string
	if err != nil || task.Id == 0 {
		logger.Errorf("编辑任务#获取任务详情失败#任务ID-%d", id)
		result = jsonResp.Success(utils.SuccessContent, nil)
	} else {
		result = jsonResp.Success(utils.SuccessContent, task)
	}
	c.String(http.StatusOK, result)
}

// 保存任务
func Store(c *gin.Context) {
	var form TaskForm
	if err := c.ShouldBind(&form); err != nil {
		base.RespondValidationError(c, err)
		return
	}

	taskModel := models.Task{}
	var id = form.Id
	nameExists, err := taskModel.NameExist(form.Name, form.Id)
	if err != nil {
		base.RespondErrorWithDefaultMsg(c, err)
		return
	}
	if nameExists {
		base.RespondError(c, i18n.T(c, "task_name_exists"))
		return
	}

	if form.Protocol == models.TaskRPC && form.HostId == "" {
		base.RespondError(c, i18n.T(c, "select_hostname"))
		return
	}

	taskModel.Name = form.Name
	taskModel.Protocol = form.Protocol
	// 清理命令中的 HTML 实体编码
	originalCmd := strings.TrimSpace(form.Command)
	cleanedCmd := utils.CleanHTMLEntities(originalCmd)
	if originalCmd != cleanedCmd {
		logger.Infof("[HTML Entity Cleaned] Task: %s, Original length: %d, Cleaned length: %d", form.Name, len(originalCmd), len(cleanedCmd))
	}
	taskModel.Command = cleanedCmd
	taskModel.Timeout = form.Timeout
	taskModel.Tag = form.Tag
	taskModel.Remark = form.Remark
	taskModel.Multi = form.Multi
	taskModel.RetryTimes = form.RetryTimes
	taskModel.RetryInterval = form.RetryInterval
	taskModel.NotifyStatus = form.NotifyStatus
	notifyTypeMask, err := models.NormalizeNotifyTypeMask(form.NotifyType)
	if err != nil {
		base.RespondError(c, i18n.T(c, "param_error"))
		return
	}
	taskModel.NotifyType = notifyTypeMask
	taskModel.NotifyReceiverId = form.NotifyReceiverId
	taskModel.NotifyKeyword = form.NotifyKeyword
	taskModel.Spec = form.Spec
	taskModel.Level = form.Level
	taskModel.DependencyStatus = form.DependencyStatus
	taskModel.DependencyTaskId = strings.TrimSpace(form.DependencyTaskId)
	if taskModel.NotifyStatus > 0 {
		receiverId := strings.TrimSpace(taskModel.NotifyReceiverId)
		if taskModel.NotifyType&(models.NotifyTypeMailMask|models.NotifyTypeSlackMask) != 0 && receiverId == "" {
			base.RespondError(c, i18n.T(c, "select_at_least_one_receiver"))
			return
		}
	}
	taskModel.HttpMethod = form.HttpMethod
	if taskModel.Protocol == models.TaskHTTP {
		command := strings.ToLower(taskModel.Command)
		if !strings.HasPrefix(command, "http://") && !strings.HasPrefix(command, "https://") {
			base.RespondError(c, i18n.T(c, "invalid_url"))
			return
		}
		if taskModel.Timeout > 300 {
			base.RespondError(c, i18n.T(c, "http_task_timeout_max_300"))
			return
		}
	}

	if taskModel.RetryTimes > 10 || taskModel.RetryTimes < 0 {
		base.RespondError(c, i18n.T(c, "retry_times_range_0_10"))
		return
	}

	if taskModel.RetryInterval > 3600 || taskModel.RetryInterval < 0 {
		base.RespondError(c, i18n.T(c, "retry_interval_range_0_3600"))
		return
	}

	if taskModel.DependencyStatus != models.TaskDependencyStatusStrong &&
		taskModel.DependencyStatus != models.TaskDependencyStatusWeak {
		base.RespondError(c, i18n.T(c, "select_dependency"))
		return
	}

	if taskModel.Level == models.TaskLevelParent {
		err = utils.PanicToError(func() {
			cron.Parse(form.Spec)
		})
		if err != nil {
			base.RespondError(c, i18n.T(c, "crontab_parse_failed"), err)
			return
		}
	} else {
		taskModel.DependencyTaskId = ""
		taskModel.Spec = ""
	}

	if id > 0 && taskModel.DependencyTaskId != "" {
		dependencyTaskIds := strings.Split(taskModel.DependencyTaskId, ",")
		if utils.InStringSlice(dependencyTaskIds, strconv.Itoa(id)) {
			base.RespondError(c, i18n.T(c, "cannot_set_self_as_child"))
			return
		}
	}

	if id == 0 {
		taskModel.Status = models.Running
		logger.Infof("[Task Create] Before Create - Multi: %d", taskModel.Multi)
		id, err = taskModel.Create()
		if err == nil {
			// 立即读取验证
			verifyTask, _ := taskModel.Detail(id)
			logger.Infof("[Task Create] After Create - ID: %d, Multi in DB: %d", id, verifyTask.Multi)
		}
	} else {
		logger.Infof("[Task Update] Before Update - ID: %d, Multi: %d", id, taskModel.Multi)
		_, err = taskModel.UpdateBean(id)
		if err == nil {
			// 立即读取验证
			verifyTask, _ := taskModel.Detail(id)
			logger.Infof("[Task Update] After Update - ID: %d, Multi in DB: %d", id, verifyTask.Multi)
		}
	}

	if err != nil {
		base.RespondError(c, i18n.T(c, "save_failed"), err)
		return
	}

	taskHostModel := new(models.TaskHost)
	if form.Protocol == models.TaskRPC {
		hostIdStrList := strings.Split(form.HostId, ",")
		hostIds := make([]int, len(hostIdStrList))
		for i, hostIdStr := range hostIdStrList {
			hostIds[i], _ = strconv.Atoi(hostIdStr)
		}
		_ = taskHostModel.Add(id, hostIds)
	} else {
		_ = taskHostModel.Remove(id)
	}

	status, _ := taskModel.GetStatus(id)
	if status == models.Enabled && taskModel.Level == models.TaskLevelParent {
		addTaskToTimer(id)
	}

	base.RespondSuccess(c, i18n.T(c, "save_success"), nil)
}

// 删除任务
func Remove(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	taskModel := new(models.Task)
	_, err := taskModel.Delete(id)
	if err != nil {
		base.RespondErrorWithDefaultMsg(c, err)
	} else {
		taskHostModel := new(models.TaskHost)
		_ = taskHostModel.Remove(id)
		service.ServiceTask.Remove(id)
		base.RespondSuccessWithDefaultMsg(c, nil)
	}
}

// 激活任务
func Enable(c *gin.Context) {
	changeStatus(c, models.Enabled)
}

// 暂停任务
func Disable(c *gin.Context) {
	changeStatus(c, models.Disabled)
}

// 手动运行任务
func Run(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	taskModel := new(models.Task)
	task, err := taskModel.Detail(id)
	if err != nil || task.Id <= 0 {
		base.RespondError(c, i18n.T(c, "get_task_detail_failed"), err)
	} else {
		task.Spec = i18n.T(c, "manual_run")
		service.ServiceTask.Run(task)
		base.RespondSuccess(c, i18n.T(c, "task_started_check_log"), nil)
	}
}

// 批量启用任务
func BatchEnable(c *gin.Context) {
	batchChangeStatus(c, models.Enabled)
}

// 批量禁用任务
func BatchDisable(c *gin.Context) {
	batchChangeStatus(c, models.Disabled)
}

// 批量改变任务状态
func batchChangeStatus(c *gin.Context, status models.Status) {
	var form struct {
		Ids []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		base.RespondError(c, i18n.T(c, "param_error"))
		return
	}

	taskModel := new(models.Task)
	successCount := 0
	for _, id := range form.Ids {
		_, err := taskModel.Update(id, models.CommonMap{
			"status": status,
		})
		if err == nil {
			successCount++
			if status == models.Enabled {
				addTaskToTimer(id)
			} else {
				service.ServiceTask.Remove(id)
			}
		}
	}

	base.RespondSuccess(c, i18n.T(c, "operation_success"), map[string]interface{}{
		"success_count": successCount,
		"total_count":   len(form.Ids),
	})
}

// 批量删除任务
func BatchRemove(c *gin.Context) {
	var form struct {
		Ids []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&form); err != nil {
		base.RespondError(c, i18n.T(c, "param_error"))
		return
	}

	taskModel := new(models.Task)
	taskHostModel := new(models.TaskHost)
	successCount := 0
	for _, id := range form.Ids {
		_, err := taskModel.Delete(id)
		if err == nil {
			successCount++
			_ = taskHostModel.Remove(id)
			service.ServiceTask.Remove(id)
		}
	}

	base.RespondSuccess(c, "操作成功", map[string]interface{}{
		"success_count": successCount,
		"total_count":   len(form.Ids),
	})
}

// 改变任务状态
func changeStatus(c *gin.Context, status models.Status) {
	id, _ := strconv.Atoi(c.Param("id"))
	taskModel := new(models.Task)
	_, err := taskModel.Update(id, models.CommonMap{
		"status": status,
	})
	if err != nil {
		base.RespondErrorWithDefaultMsg(c, err)
	} else {
		if status == models.Enabled {
			addTaskToTimer(id)
		} else {
			service.ServiceTask.Remove(id)
		}
		base.RespondSuccessWithDefaultMsg(c, nil)
	}
}

// 添加任务到定时器
func addTaskToTimer(id int) {
	taskModel := new(models.Task)
	task, err := taskModel.Detail(id)
	if err != nil {
		logger.Error(err)
		return
	}

	service.ServiceTask.RemoveAndAdd(task)
}

// 解析查询参数
func parseQueryParams(c *gin.Context) models.CommonMap {
	var params models.CommonMap = models.CommonMap{}
	id, _ := strconv.Atoi(c.Query("id"))
	hostId, _ := strconv.Atoi(c.Query("host_id"))
	protocol, _ := strconv.Atoi(c.Query("protocol"))
	status, _ := strconv.Atoi(c.Query("status"))
	params["Id"] = id
	params["HostId"] = hostId
	params["Name"] = strings.TrimSpace(c.Query("name"))
	params["Protocol"] = protocol
	params["Tag"] = strings.TrimSpace(c.Query("tag"))
	if status >= 0 {
		status -= 1
	}
	params["Status"] = status
	base.ParsePageAndPageSize(c, params)

	return params
}
