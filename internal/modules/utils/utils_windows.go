//go:build windows
// +build windows

package utils

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"golang.org/x/net/context"
)

type Result struct {
	output string
	err    error
}

// 执行shell命令，可设置执行超时时间
func ExecShell(ctx context.Context, command string) (string, error) {
	// 清理可能存在的 HTML 实体编码,防止 &quot; 等导致命令执行失败
	// 例如: del &quot;C:\file.txt&quot; -> del "C:\file.txt"
	command = CleanHTMLEntities(command)

	// 使用 cmd.exe，通过 CmdLine 直接传递完整命令行，绕过 Go 的参数转义
	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		CmdLine:    `cmd /c "` + command + `"`,
	}
	// 设置工作目录为用户家目录
	if homeDir, err := os.UserHomeDir(); err == nil {
		cmd.Dir = homeDir
	} else {
		cmd.Dir = os.TempDir()
	}
	var resultChan chan Result = make(chan Result)
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- Result{string(output), err}
	}()
	select {
	case <-ctx.Done():
		if cmd.Process.Pid > 0 {
			exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
			cmd.Process.Kill()
		}
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		return ConvertEncoding(result.output), result.err
	}
}

func ConvertEncoding(outputGBK string) string {
	// windows平台编码为gbk，需转换为utf8才能入库
	outputUTF8, ok := GBK2UTF8(outputGBK)
	if ok {
		return outputUTF8
	}

	return outputGBK
}
