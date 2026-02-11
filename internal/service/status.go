package service

type SchedulerStatus struct {
	Running              bool `json:"running"`
	EntryCount           int  `json:"entryCount"`
	ConcurrencyUsed      int  `json:"concurrencyUsed"`
	ConcurrencyCap       int  `json:"concurrencyCap"`
	RunningInstanceCount int  `json:"runningInstanceCount"`
}

func GetSchedulerStatus() SchedulerStatus {
	status := SchedulerStatus{}
	if serviceCron == nil {
		return status
	}
	status.Running = true
	status.EntryCount = len(serviceCron.Entries())
	if concurrencyQueue.queue != nil {
		status.ConcurrencyUsed = len(concurrencyQueue.queue)
		status.ConcurrencyCap = cap(concurrencyQueue.queue)
	}
	runningCount := 0
	runInstance.m.Range(func(_, _ interface{}) bool {
		runningCount++
		return true
	})
	status.RunningInstanceCount = runningCount
	return status
}

