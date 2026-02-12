//go:build !windows
// +build !windows

package utils

import (
	"context"
	"strings"
	"testing"
	"time"
)

// 集成测试：模拟真实场景 - 任务超时但需要看到已执行的输出
func TestExecShell_RealWorldScenario_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 模拟一个数据库备份任务：先输出开始信息，然后执行耗时操作
	command := `
		echo "=== Database Backup Started ==="
		echo "Connecting to database..."
		echo "Dumping table: users"
		echo "Dumping table: orders"
		sleep 5
		echo "Backup completed"
	`

	output, err := ExecShell(ctx, command)

	if err == nil {
		t.Fatal("Expected timeout error")
	}

	// 验证能看到超时前的所有输出
	requiredOutputs := []string{
		"Database Backup Started",
		"Connecting to database",
		"Dumping table: users",
		"Dumping table: orders",
	}

	for _, expected := range requiredOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Missing expected output: %s\nGot: %s", expected, output)
		}
	}

	// 不应该包含超时后的输出
	if strings.Contains(output, "Backup completed") {
		t.Error("Should not contain output after timeout")
	}

	t.Logf("✓ Successfully captured partial output on timeout:\n%s", output)
}

// 集成测试：手动停止任务
func TestExecShell_RealWorldScenario_ManualStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	command := `
		echo "Task started"
		for i in {1..20}; do
			echo "Processing record $i"
			sleep 0.2
		done
		echo "Task finished"
	`

	resultChan := make(chan struct {
		output string
		err    error
	})

	go func() {
		output, err := ExecShell(ctx, command)
		resultChan <- struct {
			output string
			err    error
		}{output, err}
	}()

	// 模拟用户在 1 秒后点击停止按钮
	time.Sleep(1 * time.Second)
	cancel()

	result := <-resultChan

	if result.err == nil {
		t.Fatal("Expected cancellation error")
	}

	// 应该能看到部分处理记录
	if !strings.Contains(result.output, "Task started") {
		t.Error("Missing 'Task started'")
	}

	if !strings.Contains(result.output, "Processing record") {
		t.Error("Missing processing records")
	}

	recordCount := strings.Count(result.output, "Processing record")
	if recordCount < 3 {
		t.Errorf("Expected at least 3 records processed, got %d", recordCount)
	}

	t.Logf("✓ Captured %d records before manual stop:\n%s", recordCount, result.output)
}
