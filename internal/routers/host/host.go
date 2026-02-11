package host

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tabortao/gocron/internal/models"
	"github.com/tabortao/gocron/internal/modules/i18n"
	"github.com/tabortao/gocron/internal/modules/logger"
	"github.com/tabortao/gocron/internal/modules/rpc/client"
	"github.com/tabortao/gocron/internal/modules/rpc/grpcpool"
	rpc "github.com/tabortao/gocron/internal/modules/rpc/proto"
	"github.com/tabortao/gocron/internal/modules/utils"
	"github.com/tabortao/gocron/internal/routers/base"
	"github.com/tabortao/gocron/internal/service"
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

	base.RespondSuccess(c, utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  hosts,
	})
}

// All 获取所有主机
func All(c *gin.Context) {
	hostModel := new(models.Host)
	hostModel.PageSize = -1
	hosts, err := hostModel.List(models.CommonMap{})
	if err != nil {
		logger.Error(err)
	}

	base.RespondSuccess(c, utils.SuccessContent, hosts)
}

// Detail 主机详情
func Detail(c *gin.Context) {
	hostModel := new(models.Host)
	id, _ := strconv.Atoi(c.Param("id"))
	err := hostModel.Find(id)
	if err != nil || hostModel.Id == 0 {
		logger.Errorf("获取主机详情失败#主机id-%d", id)
		base.RespondSuccess(c, utils.SuccessContent, nil)
	} else {
		base.RespondSuccess(c, utils.SuccessContent, hostModel)
	}
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
		base.RespondValidationError(c, err)
		return
	}

	hostModel := new(models.Host)
	id := form.Id
	nameExist, err := hostModel.NameExists(form.Name, form.Id)
	if err != nil {
		base.RespondError(c, i18n.T(c, "operation_failed"), err)
		return
	}
	if nameExist {
		base.RespondError(c, i18n.T(c, "hostname_exists"))
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
			base.RespondError(c, i18n.T(c, "host_not_exist"))
			return
		}
		_, err = hostModel.UpdateBean(id)
	} else {
		isCreate = true
		id, err = hostModel.Create()
	}
	if err != nil {
		base.RespondError(c, i18n.T(c, "save_failed"), err)
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
			base.RespondError(c, i18n.T(c, "refresh_task_host_failed"), err)
			return
		}
		service.ServiceTask.BatchAdd(tasks)
	}

	base.RespondSuccess(c, i18n.T(c, "save_success"), nil)
}

// Remove 删除主机
func Remove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		base.RespondError(c, i18n.T(c, "param_error"), err)
		return
	}
	taskHostModel := new(models.TaskHost)
	exist, err := taskHostModel.HostIdExist(id)
	if err != nil {
		base.RespondError(c, i18n.T(c, "operation_failed"), err)
		return
	}
	if exist {
		base.RespondError(c, i18n.T(c, "host_in_use_cannot_delete"))
		return
	}

	hostModel := new(models.Host)
	err = hostModel.Find(int(id))
	if err != nil {
		base.RespondError(c, i18n.T(c, "host_not_exist"))
		return
	}

	_, err = hostModel.Delete(id)
	if err != nil {
		base.RespondError(c, i18n.T(c, "operation_failed"), err)
		return
	}

	addr := fmt.Sprintf("%s:%d", hostModel.Name, hostModel.Port)
	grpcpool.Pool.Release(addr)

	base.RespondSuccess(c, i18n.T(c, "operation_success"), nil)
}

// Ping 测试主机是否可连接
func Ping(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	hostModel := new(models.Host)
	err := hostModel.Find(id)
	if err != nil || hostModel.Id <= 0 {
		base.RespondError(c, i18n.T(c, "host_not_exist"), err)
		return
	}

	taskReq := &rpc.TaskRequest{}
	taskReq.Command = testConnectionCommand
	taskReq.Timeout = testConnectionTimeout
	output, err := client.Exec(hostModel.Name, hostModel.Port, taskReq)
	if err != nil {
		var loopbackErr error
		if hostModel.Name != "127.0.0.1" && (shouldTryLoopback(hostModel.Name) || isSameAsRequestHost(c, hostModel.Name)) {
			_, loopbackErr = client.Exec("127.0.0.1", hostModel.Port, taskReq)
			if loopbackErr == nil {
				oldName := hostModel.Name
				hostModel.Name = "127.0.0.1"
				_, _ = hostModel.UpdateBean(hostModel.Id)
				grpcpool.Pool.Release(fmt.Sprintf("%s:%d", oldName, hostModel.Port))

				if i18n.GetLocale(c) == i18n.EnUS {
					base.RespondSuccess(c, "Connection successful (auto-fixed host to 127.0.0.1)", nil)
				} else {
					base.RespondSuccess(c, "连接成功（已自动将主机名修正为 127.0.0.1）", nil)
				}
				return
			}
		}
		msg := i18n.T(c, "connection_failed") + "-" + err.Error() + " " + output
		if isConnRefused(err) {
			if i18n.GetLocale(c) == i18n.EnUS {
				msg += " (Port refused: gocron-node is not listening on this address/port. Try starting gocron-node with -s 0.0.0.0:5921 and allow firewall inbound 5921.)"
			} else {
				msg += "（端口被拒绝：通常是 gocron-node 未启动，或仅监听 127.0.0.1。建议用 -s 0.0.0.0:5921 启动，并放行防火墙 5921 入站）"
			}
		}
		if loopbackErr != nil {
			if i18n.GetLocale(c) == i18n.EnUS {
				msg += " (Loopback check failed: 127.0.0.1:" + strconv.Itoa(hostModel.Port) + " -> " + loopbackErr.Error() + ")"
			} else {
				msg += "（同时检测 127.0.0.1:" + strconv.Itoa(hostModel.Port) + " 也失败：" + loopbackErr.Error() + "）"
			}
		}
		logger.Error(msg)
		base.RespondError(c, msg, err)
		return
	}
	base.RespondSuccess(c, i18n.T(c, "connection_success"), nil)
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

func shouldTryLoopback(host string) bool {
	if host == "" {
		return false
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP == nil {
			continue
		}
		if ipNet.IP.Equal(ip) {
			return true
		}
	}
	if ip.IsLoopback() {
		return true
	}
	return false
}

func isSameAsRequestHost(c *gin.Context, host string) bool {
	if c == nil || c.Request == nil {
		return false
	}
	reqHost := strings.TrimSpace(c.Request.Host)
	if reqHost == "" {
		return false
	}
	if h, _, err := net.SplitHostPort(reqHost); err == nil && h != "" {
		reqHost = h
	}
	return strings.EqualFold(strings.TrimSpace(host), reqHost)
}

func isConnRefused(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "actively refused") || strings.Contains(s, "connection refused")
}
