package service

// https://github.com/gocronx-team/gocron/issues/66

import (
	"sync"
	"testing"
	"time"

	"github.com/gocronx-team/gocron/internal/models"
)

// TestIssue66RaceCondition 重现 Issue #66: 任务无法单实例运行
// 问题：beforeExecJob检查和createJob添加实例标记之间存在竞态条件
func TestIssue66RaceCondition(t *testing.T) {
	t.Run("旧实现-重现竞态条件bug", func(t *testing.T) {
		runInstance = Instance{}

		task := models.Task{
			Id:    1,
			Name:  "测试任务",
			Multi: 0, // 禁止并发
		}

		var executedCount int
		var mu sync.Mutex
		var wg sync.WaitGroup

		// 模拟快速点击两次"手动执行" - 使用旧的实现方式
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				// 旧实现：分离的检查和添加
				// 1. beforeExecJob 检查
				taskLogId := int64(0)
				if task.Multi == 0 && runInstance.has(task.Id) {
					t.Logf("执行%d: 任务已在运行中，取消本次执行", index)
					return
				}
				taskLogId = int64(index + 1)

				// ⚠️ 竞态条件窗口：在检查和添加之间
				time.Sleep(1 * time.Millisecond)

				// 2. 添加实例标记
				if task.Multi == 0 {
					runInstance.add(task.Id)
					defer runInstance.done(task.Id)
				}

				// 3. 执行任务
				mu.Lock()
				executedCount++
				mu.Unlock()

				t.Logf("执行%d: 任务开始执行, taskLogId=%d", index, taskLogId)
				time.Sleep(10 * time.Millisecond)
				t.Logf("执行%d: 任务执行完成", index)
			}(i)
		}

		wg.Wait()

		if executedCount > 1 {
			t.Logf("✅ 成功重现Bug：期望执行1次，实际执行了%d次", executedCount)
		}
	})

	t.Run("新实现-修复竞态条件", func(t *testing.T) {
		runInstance = Instance{}

		task := models.Task{
			Id:    2,
			Name:  "测试任务",
			Multi: 0,
		}

		var executedCount int
		var canceledCount int
		var mu sync.Mutex
		var wg sync.WaitGroup

		// 模拟快速点击两次"手动执行" - 使用新的原子实现
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				// 新实现：原子的检查和添加
				if task.Multi == 0 {
					if !runInstance.tryAdd(task.Id) {
						mu.Lock()
						canceledCount++
						mu.Unlock()
						t.Logf("执行%d: 任务已在运行中，取消本次执行", index)
						return
					}
					defer runInstance.done(task.Id)
				}

				// 执行任务
				mu.Lock()
				executedCount++
				mu.Unlock()

				t.Logf("执行%d: 任务开始执行", index)
				time.Sleep(10 * time.Millisecond)
				t.Logf("执行%d: 任务执行完成", index)
			}(i)
		}

		wg.Wait()

		if executedCount == 1 && canceledCount == 1 {
			t.Logf("✅ Bug已修复！执行%d次，取消%d次", executedCount, canceledCount)
		} else {
			t.Errorf("❌ 修复失败！执行%d次，取消%d次（期望：执行1次，取消1次）", executedCount, canceledCount)
		}
	})
}
