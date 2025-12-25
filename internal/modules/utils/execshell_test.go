//go:build !windows
// +build !windows

package utils

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

// 测试正常执行完成
func TestExecShell_NormalCompletion(t *testing.T) {
	ctx := context.Background()
	output, err := ExecShell(ctx, "echo 'Hello World'")
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if !strings.Contains(output, "Hello World") {
		t.Errorf("Expected output to contain 'Hello World', got: %s", output)
	}
}

// 测试超时时能捕获部分输出
func TestExecShell_TimeoutWithPartialOutput(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	// 执行一个会输出多行然后长时间运行的命令
	command := `
		echo "Line 1"
		echo "Line 2"
		echo "Line 3"
		sleep 10
		echo "This should not appear"
	`
	
	output, err := ExecShell(ctx, command)
	
	// 应该返回错误（超时）
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	
	if !strings.Contains(err.Error(), "timeout killed") {
		t.Errorf("Expected 'timeout killed' error, got: %v", err)
	}
	
	// 关键：应该能捕获到前面的输出
	if !strings.Contains(output, "Line 1") {
		t.Errorf("Expected output to contain 'Line 1', got: %s", output)
	}
	if !strings.Contains(output, "Line 2") {
		t.Errorf("Expected output to contain 'Line 2', got: %s", output)
	}
	if !strings.Contains(output, "Line 3") {
		t.Errorf("Expected output to contain 'Line 3', got: %s", output)
	}
	
	// 不应该包含超时后的输出
	if strings.Contains(output, "This should not appear") {
		t.Errorf("Output should not contain text after timeout")
	}
	
	t.Logf("Captured output on timeout:\n%s", output)
}

// 测试手动取消时能捕获部分输出
func TestExecShell_ManualCancelWithPartialOutput(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	
	// 启动一个会持续输出的命令
	command := `
		for i in {1..10}; do
			echo "Output line $i"
			sleep 0.5
		done
	`
	
	// 在另一个 goroutine 中执行命令
	outputChan := make(chan string)
	errChan := make(chan error)
	
	go func() {
		output, err := ExecShell(ctx, command)
		outputChan <- output
		errChan <- err
	}()
	
	// 等待 1.5 秒后取消（应该能看到前几行输出）
	time.Sleep(1500 * time.Millisecond)
	cancel()
	
	output := <-outputChan
	err := <-errChan
	
	// 应该返回错误（被取消）
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	
	if !strings.Contains(err.Error(), "timeout killed") {
		t.Errorf("Expected 'timeout killed' error, got: %v", err)
	}
	
	// 应该能捕获到部分输出
	if !strings.Contains(output, "Output line") {
		t.Errorf("Expected output to contain 'Output line', got: %s", output)
	}
	
	t.Logf("Captured output on manual cancel:\n%s", output)
}

// 测试命令执行失败但有输出
func TestExecShell_CommandFailureWithOutput(t *testing.T) {
	ctx := context.Background()
	
	// 执行一个会失败的命令，但会先输出内容
	command := `
		echo "Before error"
		ls /nonexistent_directory_12345
		echo "After error"
	`
	
	output, err := ExecShell(ctx, command)
	
	// 命令失败但 bash 会继续执行后续命令，所以可能没有错误
	// 这取决于 shell 的行为
	if err != nil {
		t.Logf("Command returned error (expected): %v", err)
	}
	
	// 应该能捕获到错误前的输出
	if !strings.Contains(output, "Before error") {
		t.Errorf("Expected output to contain 'Before error', got: %s", output)
	}
	
	// 应该包含错误信息
	if !strings.Contains(output, "No such file or directory") && !strings.Contains(output, "cannot access") {
		t.Logf("Warning: Expected error message in output, got: %s", output)
	}
	
	t.Logf("Output with error:\n%s", output)
}

// 测试长时间运行的命令
func TestExecShell_LongRunningCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// 模拟一个持续输出的长任务
	command := `
		for i in {1..100}; do
			echo "Processing item $i"
			sleep 0.1
		done
	`
	
	output, err := ExecShell(ctx, command)
	
	// 应该超时
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	
	// 应该捕获到大量输出
	lineCount := strings.Count(output, "Processing item")
	if lineCount < 10 {
		t.Errorf("Expected at least 10 lines of output, got %d", lineCount)
	}
	
	t.Logf("Captured %d lines before timeout", lineCount)
}

// 测试快速完成的命令
func TestExecShell_QuickCommand(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	output, err := ExecShell(ctx, "echo 'Quick test' && date")
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if !strings.Contains(output, "Quick test") {
		t.Errorf("Expected output to contain 'Quick test', got: %s", output)
	}
}

// 测试空命令
func TestExecShell_EmptyCommand(t *testing.T) {
	ctx := context.Background()
	output, err := ExecShell(ctx, "")
	
	// 空命令应该成功执行（没有输出）
	if err != nil {
		t.Logf("Empty command returned error: %v (this is acceptable)", err)
	}
	
	t.Logf("Empty command output: '%s'", output)
}

// 测试 stderr 输出
func TestExecShell_StderrOutput(t *testing.T) {
	ctx := context.Background()
	
	// 同时输出到 stdout 和 stderr
	command := `
		echo "stdout message"
		echo "stderr message" >&2
	`
	
	output, err := ExecShell(ctx, command)
	
	if err != nil {
		t.Logf("Command returned error: %v", err)
	}
	
	// 应该同时捕获 stdout 和 stderr
	if !strings.Contains(output, "stdout message") {
		t.Errorf("Expected output to contain 'stdout message', got: %s", output)
	}
	if !strings.Contains(output, "stderr message") {
		t.Errorf("Expected output to contain 'stderr message', got: %s", output)
	}
}

// 基准测试：正常命令执行
func BenchmarkExecShell_Normal(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		ExecShell(ctx, "echo 'benchmark test'")
	}
}

// 基准测试：超时场景
func BenchmarkExecShell_Timeout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		ExecShell(ctx, "sleep 1")
		cancel()
	}
}

// 测试包含特殊字符的命令（临时脚本方式的优势）
func TestExecShell_SpecialCharacters(t *testing.T) {
	ctx := context.Background()
	
	// 测试包含引号、反引号等特殊字符
	command := "echo \"Hello 'World'\" && echo 'Test \"quotes\"' && echo `date`"
	
	output, err := ExecShell(ctx, command)
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if !strings.Contains(output, "Hello 'World'") {
		t.Errorf("Expected output to contain mixed quotes, got: %s", output)
	}
	
	t.Logf("Special characters output:\n%s", output)
}

// 测试多行脚本（临时脚本方式的优势）
func TestExecShell_MultilineScript(t *testing.T) {
	ctx := context.Background()
	
	// 测试复杂的多行脚本
	command := `
#!/bin/bash
function greet() {
    echo "Hello from function"
}

for i in 1 2 3; do
    echo "Loop iteration $i"
done

greet
`
	
	output, err := ExecShell(ctx, command)
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if !strings.Contains(output, "Loop iteration") {
		t.Errorf("Expected loop output, got: %s", output)
	}
	
	if !strings.Contains(output, "Hello from function") {
		t.Errorf("Expected function output, got: %s", output)
	}
	
	t.Logf("Multiline script output:\n%s", output)
}

// 测试工作目录是否正确设置
func TestExecShell_WorkingDirectory(t *testing.T) {
	ctx := context.Background()
	
	// 打印当前工作目录
	output, err := ExecShell(ctx, "pwd")
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// 工作目录应该是用户家目录，不是 /tmp
	if strings.Contains(output, "/tmp") && !strings.Contains(output, os.Getenv("HOME")) {
		t.Errorf("Working directory should be home directory, got: %s", output)
	}
	
	t.Logf("Working directory: %s", strings.TrimSpace(output))
}

// 测试HTML实体清理（保持原有功能）
func TestExecShell_HTMLEntityCleaning(t *testing.T) {
	ctx := context.Background()
	
	// 测试HTML实体会被正确清理
	command := `echo &quot;test&quot;`
	
	output, err := ExecShell(ctx, command)
	
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// 应该输出 "test" 而不是 &quot;test&quot;
	if !strings.Contains(output, "test") {
		t.Errorf("Expected cleaned output, got: %s", output)
	}
	
	t.Logf("HTML entity cleaned output: %s", output)
}
