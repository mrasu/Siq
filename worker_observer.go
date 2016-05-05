package siq

import (
	"sync"
	"github.com/mrasu/Siq/workers"
)

type WorkerObserver struct {
	workingWorker map[*workers.Worker]struct{}
	workers       []*workers.Worker
	m             sync.Locker
	ss            *SystemSiq
}

func newWorkerObserver(s *Siq) *WorkerObserver {
	wo := &WorkerObserver{}
	wo.init(s)
	return wo
}

func (wo *WorkerObserver) init(s *Siq) {
	wo.workingWorker = map[*workers.Worker]struct{}{}
	wo.workers = []*workers.Worker{}
	wo.m = &sync.Mutex{}
	wo.ss = NewSystemSiq(wo, s)
}

func (wo *WorkerObserver) Add(w *workers.Worker) {
	wo.m.Lock()
	defer wo.m.Unlock()
	wo.workers = append(wo.workers, w)
}

func (wo *WorkerObserver) Delete(tw *workers.Worker) {
	wo.m.Lock()
	defer wo.m.Unlock()
	for i, w := range(wo.workers) {
		if w == tw {
			wo.workers = append(wo.workers[:i], wo.workers[i+1:]...)
			break
		}
	}
}

func (wo *WorkerObserver) Count() int {
	return len(wo.workers)
}

func (wo *WorkerObserver) LockIdleWorker(start int) (int, *workers.Worker){
	wo.m.Lock()
	defer wo.m.Unlock()

	if len(wo.workers) == 0 {
		return -1, nil
	}

	if len(wo.workers) <= start {
		start = 0
	}
	next := start
	for true {
		w := wo.workers[next]
		if _, working := wo.workingWorker[w]; !working {
			wo.workingWorker[w] = struct{}{}
			return next, w
		}

		next += 1
		if len(wo.workers) - 1 < next {
			next = 0
		}
		if next == start {
			break
		}
	}
	return -1, nil
}

func (wo *WorkerObserver) Unlock(w *workers.Worker) {
	wo.m.Lock()
	defer wo.m.Unlock()
	if _, ok := wo.workingWorker[w]; ok {
		delete(wo.workingWorker, w)
	}
}

func (wo *WorkerObserver) NotifyFail(w *workers.Worker) {
	wo.Delete(w)
	wo.ss.AddDyingWorker(w)
}