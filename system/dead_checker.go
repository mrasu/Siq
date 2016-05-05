package system

import (
	"time"
	"sync"
	"github.com/mrasu/Siq/workers"
	"fmt"
)

type DeadChecker struct {
	PollCount       int
	AliveNotifyFunc func(*workers.Worker)
	DeadNotifyFunc  func(*workers.Worker)

	dyingWorkers    map[*workers.Worker]*failInfo
	dyingWorkerLock sync.Locker

	isPolling       bool
}

type failInfo struct {
	next  time.Time
	count int
}

func NewDeadChecker(pollCount int, aliveNoticeFunc func(*workers.Worker), deadNoticeFunc func(*workers.Worker)) *DeadChecker{
	return &DeadChecker{
		PollCount: pollCount,
		AliveNotifyFunc: aliveNoticeFunc,
		DeadNotifyFunc: deadNoticeFunc,

		dyingWorkers: map[*workers.Worker]*failInfo{},
		dyingWorkerLock: &sync.Mutex{},

		isPolling: false,
	}
}

func (dc *DeadChecker) AddDyingWorkers(w *workers.Worker) {
	dc.dyingWorkerLock.Lock()
	defer dc.dyingWorkerLock.Unlock()

	if 	_, exists := dc.dyingWorkers[w]; exists{
		return
	}

	dc.dyingWorkers[w] = &failInfo{next: time.Now().Add(1 * time.Second), count: 0}

	if !dc.isPolling {
		dc.pollRetry()
	}
}

func(dc *DeadChecker) removeDyingWorkers(w *workers.Worker) {
	dc.dyingWorkerLock.Lock()
	defer dc.dyingWorkerLock.Unlock()

	delete(dc.dyingWorkers, w)
}

func(dc *DeadChecker) pollRetry() {
	if dc.isPolling {
		return
	}

	go func() {
		dc.isPolling = true
		defer func() {
			dc.isPolling = false
		}()

		for true {
			workers := dc.GetWorkers()
			for _, w := range (workers) {
				if (*w).Ping() {
					fmt.Printf("Worker(%d) is alive\n", (*w).GetId())
					go dc.AliveNotifyFunc(w)
					dc.removeDyingWorkers(w)
					continue
				}

				info := dc.dyingWorkers[w]
				if info == nil {
					continue
				}

				if info.count > dc.PollCount {
					go dc.DeadNotifyFunc(w)
					dc.removeDyingWorkers(w)
					continue
				}
				fmt.Printf("Worker(%d) is not detected(count=%d)\n", (*w).GetId(), info.count)
				info.count += 1
			}
			if len(dc.dyingWorkers) == 0 {
				fmt.Println("No DyingWorker exists")
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

func(dc *DeadChecker) GetWorkers() []*workers.Worker{
	dc.dyingWorkerLock.Lock()
	defer dc.dyingWorkerLock.Unlock()

	ws := make([]*workers.Worker, 0, len(dc.dyingWorkers))
	for w := range(dc.dyingWorkers) {
		ws = append(ws, w)
	}

	return ws
}
