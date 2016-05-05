package siq

import (
	"errors"
	"fmt"
	"math/rand"
)

func StartSiq() []*Siq{
	so := newObserver(2)
	so.start()

	return so.siqList
}

type SiqObserver struct {
	siqList []*Siq
	topicMap map[string]*Siq
	backupMap map[string]*Siq
}

func newObserver(siqCount int) *SiqObserver {
	so := &SiqObserver{
		siqList: []*Siq{},
		topicMap: map[string]*Siq{},
		backupMap: map[string]*Siq{},
	}

	for i :=0; i < siqCount; i++ {
		so.siqList = append(so.siqList, newSiq(so))
	}

	return so
}

func(so *SiqObserver) start() {
	for _, s := range(so.siqList) {
		s.Start()
	}
}

func(so *SiqObserver) GetSiqByTopic(t string) *Siq {
	return so.topicMap[t]
}

func(so *SiqObserver) GetBackupSiqByTopic(t string) *Siq {
	s, ok := so.backupMap[t]

	if ok {
		return s
	}

	s = so.electBackupSiq(t)
	return s
}

func(so *SiqObserver) SetSiqTopic(t string, s *Siq) error {
	if _, exists := so.topicMap[t]; exists {
		return errors.New(fmt.Sprintf("topic already exists: %s", t))
	}
	so.topicMap[t] = s
	so.electBackupSiq(t)

	return nil
}

func (so *SiqObserver) electBackupSiq(t string) *Siq {
	if len(so.siqList) < 2 {
		return nil
	}

	mainSiq, ok := so.topicMap[t]
	if ok == false {
		return nil
	}
	i := rand.Intn(len(so.siqList))
	backupSiq := so.siqList[i]

	if mainSiq == backupSiq {
		backupSiq = so.siqList[(i + 1) % len(so.siqList)]
	}

	so.backupMap[t] = backupSiq
	return backupSiq
}