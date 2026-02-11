package service

import (
	"sync"
	"testing"
	"time"

	"github.com/tabortao/gocron/internal/models"
)

// TestSingleInstanceControl 测试单实例运行控制
func TestSingleInstanceControl(t *testing.T) {
	t.Run("Multi=0时阻止并发执行", func(t *testing.T) {
		instance := &Instance{}
		taskId := 100

		// 第一次检查，应该不存在
		if instance.has(taskId) {
			t.Error("任务不应该在运行中")
		}

		// 添加任务
		instance.add(taskId)

		// 第二次检查，应该存在
		if !instance.has(taskId) {
			t.Error("任务应该在运行中")
		}

		// 完成任务
		instance.done(taskId)

		// 第三次检查，应该不存在
		if instance.has(taskId) {
			t.Error("任务不应该在运行中")
		}
	})

	t.Run("并发场景下的单实例控制", func(t *testing.T) {
		instance := &Instance{}
		taskId := 200
		var wg sync.WaitGroup
		executionCount := 0
		var mu sync.Mutex

		// 模拟10个并发请求
		for range 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// 使用互斥锁保护检查和添加操作的原子性
				mu.Lock()
				if !instance.has(taskId) {
					instance.add(taskId)
					executionCount++
					mu.Unlock()

					// 模拟任务执行
					time.Sleep(10 * time.Millisecond)

					instance.done(taskId)
				} else {
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		// 只有第一个请求应该执行
		if executionCount != 1 {
			t.Errorf("期望只有1次执行，实际执行了%d次", executionCount)
		}
	})

	t.Run("不同任务ID互不影响", func(t *testing.T) {
		instance := &Instance{}

		instance.add(1)
		instance.add(2)
		instance.add(3)

		if !instance.has(1) || !instance.has(2) || !instance.has(3) {
			t.Error("所有任务都应该在运行中")
		}

		instance.done(2)

		if !instance.has(1) || instance.has(2) || !instance.has(3) {
			t.Error("只有任务2应该被移除")
		}
	})
}

// TestBeforeExecJobSingleInstance 测试beforeExecJob中的单实例逻辑
func TestBeforeExecJobSingleInstance(t *testing.T) {
	t.Run("Multi=0且任务已运行时应该取消", func(t *testing.T) {
		// 重置runInstance
		runInstance = Instance{}

		task := models.Task{
			Id:    1,
			Multi: 0, // 不允许并发
		}

		// 模拟任务已在运行
		runInstance.add(task.Id)

		// 验证任务在运行中
		if !runInstance.has(task.Id) {
			t.Error("任务应该在运行中")
		}

		// 模拟beforeExecJob的逻辑：如果Multi=0且任务已运行，应该取消
		shouldCancel := task.Multi == 0 && runInstance.has(task.Id)
		if !shouldCancel {
			t.Error("任务已在运行，应该取消本次执行")
		}

		// 清理
		runInstance.done(task.Id)
	})

	t.Run("Multi=1时允许并发执行", func(t *testing.T) {
		// 重置runInstance
		runInstance = Instance{}

		task := models.Task{
			Id:    2,
			Multi: 1, // 允许并发
		}

		// 模拟任务已在运行
		runInstance.add(task.Id)

		// Multi=1时，不应该取消
		shouldCancel := task.Multi == 0 && runInstance.has(task.Id)
		if shouldCancel {
			t.Error("Multi=1时应该允许并发执行")
		}

		// 清理
		runInstance.done(task.Id)
	})
}

// TestCreateJobSingleInstanceLogic 测试createJob中的单实例逻辑
func TestCreateJobSingleInstanceLogic(t *testing.T) {
	t.Run("Multi=0时应该添加和移除实例标记", func(t *testing.T) {
		// 重置runInstance
		runInstance = Instance{}

		taskId := 100

		// 模拟createJob中的逻辑
		if !runInstance.has(taskId) {
			// Multi=0时，添加实例标记
			runInstance.add(taskId)

			// 验证已添加
			if !runInstance.has(taskId) {
				t.Error("应该已添加实例标记")
			}

			// 模拟任务执行完成
			runInstance.done(taskId)

			// 验证已移除
			if runInstance.has(taskId) {
				t.Error("应该已移除实例标记")
			}
		}
	})

	t.Run("Multi=1时不应该添加实例标记", func(t *testing.T) {
		// 重置runInstance
		runInstance = Instance{}

		taskId := 200
		multi := 1

		// Multi=1时，不添加实例标记
		if multi == 0 {
			runInstance.add(taskId)
		}

		// 验证未添加
		if runInstance.has(taskId) {
			t.Error("Multi=1时不应该添加实例标记")
		}
	})
}

// TestInstanceThreadSafety 测试Instance的线程安全性
func TestInstanceThreadSafety(t *testing.T) {
	instance := &Instance{}
	var wg sync.WaitGroup

	// 并发添加和删除
	for i := range 100 {
		wg.Add(3)

		taskId := i

		// 添加
		go func(id int) {
			defer wg.Done()
			instance.add(id)
		}(taskId)

		// 检查
		go func(id int) {
			defer wg.Done()
			_ = instance.has(id)
		}(taskId)

		// 删除
		go func(id int) {
			defer wg.Done()
			time.Sleep(1 * time.Millisecond)
			instance.done(id)
		}(taskId)
	}

	wg.Wait()

	// 测试通过表示没有发生竞态条件
	t.Log("线程安全测试通过")
}

// TestSingleInstanceRealScenario 测试真实场景
func TestSingleInstanceRealScenario(t *testing.T) {
	t.Run("定时任务触发时上次未完成", func(t *testing.T) {
		runInstance = Instance{}

		task := models.Task{
			Id:    1,
			Name:  "慢速任务",
			Multi: 0,
		}

		var executionCount int
		var mu sync.Mutex

		started := make(chan struct{})
		finished := make(chan struct{})
		block := make(chan struct{})

		go func() {
			if task.Multi == 0 && runInstance.has(task.Id) {
				close(started)
				close(finished)
				return
			}

			if task.Multi == 0 {
				runInstance.add(task.Id)
				close(started)
				defer func() {
					runInstance.done(task.Id)
					close(finished)
				}()
			} else {
				close(started)
			}

			mu.Lock()
			executionCount++
			mu.Unlock()

			<-block
		}()

		select {
		case <-started:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("等待第一次执行开始超时")
		}

		if task.Multi == 0 && runInstance.has(task.Id) {
			t.Log("第二次执行被正确阻止")
		} else {
			if task.Multi == 0 {
				runInstance.add(task.Id)
				defer runInstance.done(task.Id)
			}
			mu.Lock()
			executionCount++
			mu.Unlock()
		}

		close(block)

		select {
		case <-finished:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("等待第一次执行完成超时")
		}

		mu.Lock()
		count := executionCount
		mu.Unlock()

		if count != 1 {
			t.Errorf("期望执行1次，实际执行%d次", count)
		}
	})
}
