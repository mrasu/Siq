package siq

import (
	"sync"
	"github.com/mrasu/Siq/workers"
	"github.com/mrasu/Siq/system"
	"fmt"
	"github.com/mrasu/Siq/surface"
)

type SystemSiq struct {
	dc *system.DeadChecker
	b  *system.Backupper
	m  sync.Locker
	s  *Siq
	wo *WorkerObserver
}

func NewSystemSiq(wo *WorkerObserver, s *Siq) *SystemSiq{
	ss := &SystemSiq{
		m: &sync.Mutex{},
		s: s,
		wo: wo,
	}

	ss.dc = system.NewDeadChecker(3, ss.fireAlive, ss.fireFail)
	ss.b = system.NewBackupper(ss.getTopic, ss.getMessages, ss.backup, ss.notifyDeadTopic)
	//ss.b.StartPolling()

	return ss
}

func(ss *SystemSiq) fireAlive(w *workers.Worker) {
	ss.wo.Add(w)
}

func(ss *SystemSiq) fireFail(w *workers.Worker) {
	fmt.Printf("Worker(%d) is dead\n", (*w).GetId())
}

func(ss *SystemSiq) AddDyingWorker(w *workers.Worker) {
	fmt.Println("AddDyingWorker")
	ss.dc.AddDyingWorkers(w)
}

func(ss *SystemSiq) getTopic() []string {
	return ss.s.getTopics()
}

func (ss *SystemSiq) getMessages(tName string, fromId int) []surface.Message {
	return ss.s.getMessages(tName, fromId)
}

func (ss *SystemSiq) backup(tName string, data []byte) error {
	targetSiq := ss.s.getBackupSiq(tName)
	err := targetSiq.updateBackup(data)

	return err
}

func (ss *SystemSiq) notifyDeadTopic(tName string) {
	fmt.Printf("Siq is dead\n")
}