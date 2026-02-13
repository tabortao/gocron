package client

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/status"

	"github.com/tabortao/gocron/internal/modules/i18n"
	"github.com/tabortao/gocron/internal/modules/logger"
	"github.com/tabortao/gocron/internal/modules/rpc/grpcpool"
	pb "github.com/tabortao/gocron/internal/modules/rpc/proto"
	"google.golang.org/grpc/codes"
)

var (
	taskCtxMap     sync.Map // 存储任务执行的 context.CancelFunc
	errUnavailable = errors.New(i18n.Translate("rpc_unavailable"))
	ErrManualStop  = errors.New("rpc_manual_stop") // 特殊错误标识，用于判断是否手动停止
)

func generateTaskUniqueKey(ip string, port int, id int64) string {
	return fmt.Sprintf("%s:%d:%d", ip, port, id)
}

func Stop(ip string, port int, id int64) {
	// 异步发送停止信号，不阻塞调用者
	go func() {
		addr := fmt.Sprintf("%s:%d", ip, port)
		c, err := grpcpool.Pool.Get(addr)
		if err != nil {
			logger.Errorf("连接服务器失败#%s#%v", addr, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err = c.Run(ctx, &pb.TaskRequest{
			Command: "__STOP__",
			Id:      id,
		})
		if err != nil {
			if status.Code(err) == codes.Unavailable {
				grpcpool.Pool.Release(addr)
			}
			logger.Errorf("发送停止信号失败#%v", err)
		}
	}()
}

func Tail(ip string, port int, id int64) (string, error) {
	return Exec(ip, port, &pb.TaskRequest{
		Command: "__TAIL__",
		Timeout: 5,
		Id:      id,
	})
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
	// RPC context: 比任务超时多5秒，给服务端时间清理进程并返回输出
	ctx, cancel := context.WithTimeout(context.Background(), timeout+5*time.Second)
	defer cancel()

	taskUniqueKey := generateTaskUniqueKey(ip, port, taskReq.Id)
	taskCtxMap.Store(taskUniqueKey, cancel)
	defer taskCtxMap.Delete(taskUniqueKey)

	resp, err := c.Run(ctx, taskReq)
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			grpcpool.Pool.Release(addr)
		}
	}

	// 处理响应：即使有错误，也要返回已产生的输出
	if err != nil {
		if resp != nil && resp.Output != "" {
			return resp.Output, parseGRPCErrorOnly(err)
		}
		return parseGRPCError(err)
	}

	if resp.Error == "" {
		return resp.Output, nil
	}

	// 检查是否是手动停止
	if resp.Error == "manual stop" {
		return resp.Output, ErrManualStop
	}

	return resp.Output, errors.New(resp.Error)
}

func parseGRPCError(err error) (string, error) {
	switch status.Code(err) {
	case codes.Unavailable:
		return "", errors.New(errUnavailable.Error() + ": " + err.Error())
	case codes.DeadlineExceeded:
		return "", errors.New(i18n.Translate("rpc_timeout"))
	case codes.Canceled:
		return "", ErrManualStop
	}
	return "", err
}

// parseGRPCErrorOnly 只返回错误，不返回输出
func parseGRPCErrorOnly(err error) error {
	switch status.Code(err) {
	case codes.Unavailable:
		return errors.New(errUnavailable.Error() + ": " + err.Error())
	case codes.DeadlineExceeded:
		return errors.New(i18n.Translate("rpc_timeout"))
	case codes.Canceled:
		return ErrManualStop
	}
	return err
}
