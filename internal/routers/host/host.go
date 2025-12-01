package host

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/models"
	"github.com/gocronx-team/gocron/internal/modules/i18n"
	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/gocronx-team/gocron/internal/modules/rpc/client"
	"github.com/gocronx-team/gocron/internal/modules/rpc/grpcpool"
	rpc "github.com/gocronx-team/gocron/internal/modules/rpc/proto"
	"github.com/gocronx-team/gocron/internal/modules/utils"
	"github.com/gocronx-team/gocron/internal/routers/base"
	"github.com/gocronx-team/gocron/internal/service"
)

const testConnectionCommand = "echo hello"
const testConnectionTimeout = 5

// Index 主机列表
func Index(c *gin.Context) {
	hostModel := new(models.Host)
	queryParams := parseQueryParams(c)
	total, err := hostModel.Total(queryParams)
	if err != nil {
		logger.Error(err)
	}
	hosts, err := hostModel.List(queryParams)
	if err != nil {
		logger.Error(err)
	}

	jsonResp := utils.JsonResponse{}
	result := jsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  hosts,
	})
	c.String(http.StatusOK, result)
}

// All 获取所有主机
func All(c *gin.Context) {
	hostModel := new(models.Host)
	hostModel.PageSize = -1
	hosts, err := hostModel.List(models.CommonMap{})
	if err != nil {
		logger.Error(err)
	}

	jsonResp := utils.JsonResponse{}
	result := jsonResp.Success(utils.SuccessContent, hosts)
	c.String(http.StatusOK, result)
}

// Detail 主机详情
func Detail(c *gin.Context) {
	hostModel := new(models.Host)
	id, _ := strconv.Atoi(c.Param("id"))
	err := hostModel.Find(id)
	jsonResp := utils.JsonResponse{}
	var result string
	if err != nil || hostModel.Id == 0 {
		logger.Errorf("获取主机详情失败#主机id-%d", id)
		result = jsonResp.Success(utils.SuccessContent, nil)
	} else {
		result = jsonResp.Success(utils.SuccessContent, hostModel)
	}
	c.String(http.StatusOK, result)
}

type HostForm struct {
	Id     int    `form:"id" json:"id"`
	Name   string `form:"name" json:"name" binding:"required,max=64"`
	Alias  string `form:"alias" json:"alias" binding:"required,max=32"`
	Port   int    `form:"port" json:"port" binding:"required,min=1,max=65535"`
	Remark string `form:"remark" json:"remark"`
}

// Store 保存、修改主机信息
func Store(c *gin.Context) {
	var form HostForm
	if err := c.ShouldBind(&form); err != nil {
		json := utils.JsonResponse{}
		result := json.CommonFailure(i18n.T(c, "form_validation_failed"))
		c.String(http.StatusOK, result)
		return
	}

	json := utils.JsonResponse{}
	hostModel := new(models.Host)
	id := form.Id
	nameExist, err := hostModel.NameExists(form.Name, form.Id)
	if err != nil {
		result := json.CommonFailure(i18n.T(c, "operation_failed"), err)
		c.String(http.StatusOK, result)
		return
	}
	if nameExist {
		result := json.CommonFailure(i18n.T(c, "hostname_exists"))
		c.String(http.StatusOK, result)
		return
	}

	hostModel.Name = strings.TrimSpace(form.Name)
	hostModel.Alias = strings.TrimSpace(form.Alias)
	hostModel.Port = form.Port
	hostModel.Remark = strings.TrimSpace(form.Remark)
	isCreate := false
	oldHostModel := new(models.Host)

	if id > 0 {
		err = oldHostModel.Find(int(id))
		if err != nil {
			result := json.CommonFailure(i18n.T(c, "host_not_exist"))
			c.String(http.StatusOK, result)
			return
		}
		_, err = hostModel.UpdateBean(id)
	} else {
		isCreate = true
		id, err = hostModel.Create()
	}
	if err != nil {
		result := json.CommonFailure(i18n.T(c, "save_failed"), err)
		c.String(http.StatusOK, result)
		return
	}

	if !isCreate {
		oldAddr := fmt.Sprintf("%s:%d", oldHostModel.Name, oldHostModel.Port)
		newAddr := fmt.Sprintf("%s:%d", hostModel.Name, hostModel.Port)
		if oldAddr != newAddr {
			grpcpool.Pool.Release(oldAddr)
		}

		taskModel := new(models.Task)
		tasks, err := taskModel.ActiveListByHostId(id)
		if err != nil {
			result := json.CommonFailure(i18n.T(c, "refresh_task_host_failed"), err)
			c.String(http.StatusOK, result)
			return
		}
		service.ServiceTask.BatchAdd(tasks)
	}

	result := json.Success(i18n.T(c, "save_success"), nil)
	c.String(http.StatusOK, result)
}

// Remove 删除主机
func Remove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	json := utils.JsonResponse{}
	var result string
	if err != nil {
		result = json.CommonFailure(i18n.T(c, "param_error"), err)
		c.String(http.StatusOK, result)
		return
	}
	taskHostModel := new(models.TaskHost)
	exist, err := taskHostModel.HostIdExist(id)
	if err != nil {
		result = json.CommonFailure(i18n.T(c, "operation_failed"), err)
		c.String(http.StatusOK, result)
		return
	}
	if exist {
		result = json.CommonFailure(i18n.T(c, "host_in_use_cannot_delete"))
		c.String(http.StatusOK, result)
		return
	}

	hostModel := new(models.Host)
	err = hostModel.Find(int(id))
	if err != nil {
		result = json.CommonFailure(i18n.T(c, "host_not_exist"))
		c.String(http.StatusOK, result)
		return
	}

	_, err = hostModel.Delete(id)
	if err != nil {
		result = json.CommonFailure(i18n.T(c, "operation_failed"), err)
		c.String(http.StatusOK, result)
		return
	}

	addr := fmt.Sprintf("%s:%d", hostModel.Name, hostModel.Port)
	grpcpool.Pool.Release(addr)

	result = json.Success(i18n.T(c, "operation_success"), nil)
	c.String(http.StatusOK, result)
}

// Ping 测试主机是否可连接
func Ping(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	hostModel := new(models.Host)
	err := hostModel.Find(id)
	json := utils.JsonResponse{}
	var result string
	if err != nil || hostModel.Id <= 0 {
		result = json.CommonFailure(i18n.T(c, "host_not_exist"), err)
		c.String(http.StatusOK, result)
		return
	}

	taskReq := &rpc.TaskRequest{}
	taskReq.Command = testConnectionCommand
	taskReq.Timeout = testConnectionTimeout
	output, err := client.Exec(hostModel.Name, hostModel.Port, taskReq)
	if err != nil {
		result = json.CommonFailure(i18n.T(c, "connection_failed")+"-"+err.Error()+" "+output, err)
	} else {
		result = json.Success(i18n.T(c, "connection_success"), nil)
	}
	c.String(http.StatusOK, result)
}

// 解析查询参数
func parseQueryParams(c *gin.Context) models.CommonMap {
	var params = models.CommonMap{}
	id, _ := strconv.Atoi(c.Query("id"))
	params["Id"] = id
	params["Name"] = strings.TrimSpace(c.Query("name"))
	base.ParsePageAndPageSize(c, params)

	return params
}
