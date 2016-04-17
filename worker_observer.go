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

func newWorkerManager() *WorkerObserver {
	wm := &WorkerObserver{}
	wm.init()
	return wm
}

func (wm *WorkerObserver) init() {
	wm.workingWorker = map[*workers.Worker]struct{}{}
	wm.workers = []*workers.Worker{}
	wm.m = &sync.Mutex{}
	wm.ss = NewSystemSiq(wm)
}

func (wm *WorkerObserver) Add(w *workers.Worker) {
	wm.m.Lock()
	defer wm.m.Unlock()
	wm.workers = append(wm.workers, w)
}

func (wm *WorkerObserver) Delete(tw *workers.Worker) {
	wm.m.Lock()
	defer wm.m.Unlock()
	for i, w := range(wm.workers) {
		if w == tw {
			wm.workers = append(wm.workers[:i], wm.workers[i+1:]...)
			break
		}
	}
}

func (wm *WorkerObserver) Count() int {
	return len(wm.workers)
}

func (wm *WorkerObserver) LockIdleWorker(start int) (int, *workers.Worker){
	wm.m.Lock()
	defer wm.m.Unlock()

	if len(wm.workers) == 0 {
		return -1, nil
	}

	if len(wm.workers) <= start {
		start = 0
	}
	next := start
	for true {
		w := wm.workers[next]
		if _, working := wm.workingWorker[w]; !working {
			wm.workingWorker[w] = struct{}{}
			return next, w
		}

		next += 1
		if len(wm.workers) - 1 < next {
			next = 0
		}
		if next == start {
			break
		}
	}
	return -1, nil
}

func (wm *WorkerObserver) Unlock(w *workers.Worker) {
	wm.m.Lock()
	defer wm.m.Unlock()
	if _, ok := wm.workingWorker[w]; ok {
		delete(wm.workingWorker, w)
	}
}

func (wm *WorkerObserver) NotifyFail(w *workers.Worker) {
	wm.Delete(w)
	wm.ss.AddDyingWorker(w)
}