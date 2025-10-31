package client

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/status"

	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/gocronx-team/gocron/internal/modules/rpc/grpcpool"
	pb "github.com/gocronx-team/gocron/internal/modules/rpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

var (
	taskMap sync.Map
)

var (
	errUnavailable = errors.New("无法连接远程服务器")
)

func generateTaskUniqueKey(ip string, port int, id int64) string {
	return fmt.Sprintf("%s:%d:%d", ip, port, id)
}

func Stop(ip string, port int, id int64) {
	key := generateTaskUniqueKey(ip, port, id)
	logger.Infof("尝试停止任务#key-%s#taskLogId-%d", key, id)
	cancel, ok := taskMap.Load(key)
	if !ok {
		logger.Warnf("未找到运行中的任务#key-%s", key)
		return
	}
	logger.Infof("找到运行中的任务，执行停止#key-%s", key)
	cancel.(context.CancelFunc)()
}

func Exec(ip string, port int, taskReq *pb.TaskRequest) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("panic#rpc/client.go:Exec#", err)
		}
	}()
	addr := fmt.Sprintf("%s:%d", ip, port)
	c, err := grpcpool.Pool.Get(addr)
	if err != nil {
		return "", err
	}
	if taskReq.Timeout <= 0 || taskReq.Timeout > 86400 {
		taskReq.Timeout = 86400
	}
	timeout := time.Duration(taskReq.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	taskUniqueKey := generateTaskUniqueKey(ip, port, taskReq.Id)
	logger.Infof("任务开始执行，存储cancel函数#key-%s#taskLogId-%d", taskUniqueKey, taskReq.Id)
	taskMap.Store(taskUniqueKey, cancel)
	defer func() {
		logger.Infof("任务执行完成，删除cancel函数#key-%s", taskUniqueKey)
		taskMap.Delete(taskUniqueKey)
	}()

	resp, err := c.Run(ctx, taskReq)
	if err != nil {
		return parseGRPCError(err)
	}

	if resp.Error == "" {
		return resp.Output, nil
	}

	return resp.Output, errors.New(resp.Error)
}

func parseGRPCError(err error) (string, error) {
	switch status.Code(err) {
	case codes.Unavailable:
		return "", errUnavailable
	case codes.DeadlineExceeded:
		return "", errors.New("执行超时, 强制结束")
	case codes.Canceled:
		return "", errors.New("手动停止")
	}
	return "", err
}
