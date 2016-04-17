package siq

import (
	"sync"
	"github.com/mrasu/Siq/workers"
	"github.com/mrasu/Siq/system"
	"fmt"
)

type SystemSiq struct {
	fc *system.DeadChecker
	m  sync.Locker
	wm             *WorkerObserver
}

func NewSystemSiq(wm *WorkerObserver) *SystemSiq{
	ss := &SystemSiq{
		m: &sync.Mutex{},
		wm: wm,
	}

	ss.fc = system.NewDeadChecker(3, ss.FireAlive, ss.FireFail)

	return ss
}

func(ss *SystemSiq) FireAlive(w *workers.Worker) {
	ss.wm.Add(w)
}

func(ss *SystemSiq) FireFail(w *workers.Worker) {

}

func(ss *SystemSiq) AddDyingWorker(w *workers.Worker) {
	fmt.Println("AddDyingWorker")
	ss.fc.AddDyingWorkers(w)
}


