package siq

import (
	"fmt"
	"time"
	"github.com/mrasu/Siq/surface"
	"github.com/mrasu/Siq/workers"
)

// Siq is distributed job queue
type Siq struct {
	topics         map[string]chan *surface.TopicCommand
	topicConsuming map[string]bool
	wo             *WorkerObserver
	so *SiqObserver
	backupData []byte
}

func newSiq(so *SiqObserver) *Siq {
	s := &Siq{
		topics: map[string]chan *surface.TopicCommand{},
		so: so,
	}
	s.wo = newWorkerObserver(s)

	return s
}

func (s *Siq) Start() {
	go func() {
		for true {
			for n, _ := range s.topics {
				s.StartConsume(n)
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

func (s *Siq) Register(w workers.Worker) {
	s.wo.Add(&w)
}

func (s *Siq) Unregister(w workers.Worker) {
	s.wo.Delete(&w)
}

func (s *Siq) Publish(tName string, m string) *Siq {
	topic := s.getTopicChannel(tName)
	if topic == nil {
		topicSiq := s.so.GetSiqByTopic(tName)
		if topicSiq == nil {
			topic = s.createTopicChannel(tName)
		} else {
			return topicSiq
		}
	}
	surface.AddTopic(topic, m)
	s.StartConsume(tName)
	return nil
}

func (s *Siq) StartConsume(tName string) {
	if _, exists := s.topicConsuming[tName]; exists {
		return
	}

	go func() {
		s.consumeAllAt(tName)
		delete(s.topicConsuming, tName)
		fmt.Printf("consumeAllAt(%s)\n", tName)
	}()
}

func (s *Siq) ConsumeTopic(tName string) string {
	t, ok := s.topics[tName]
	if !ok {
		return ""
	}

	return surface.DequeTopic(t)
}

func (s *Siq) getTopicChannel(n string) chan *surface.TopicCommand {
	if c, ok := s.topics[n]; ok {
		return c
	}
	return nil
}

func (s *Siq) createTopicChannel(n string) chan *surface.TopicCommand {
	c := surface.NewTopic(n)
	s.topics[n] = c
	err := s.so.SetSiqTopic(n, s)
	if err != nil {
		panic(err)
	}
	return c
}

func (s *Siq) Show() {
	fmt.Println("Show Topic!!")
	for _, t := range s.topics {
		surface.ShowTopic(t)
	}
}

func (s *Siq) getTopics() []string{
	siq_topics := make([]string, 0, len(s.topics))
	for t := range(s.topics) {
		siq_topics = append(siq_topics, t)
	}

	return siq_topics
}

func (s *Siq) getMessages(tName string, fromId int) []surface.Message {
	t, ok := s.topics[tName]
	if ok == false {
		return []surface.Message{}
	}

	return surface.GetTopicFrom(t, fromId)
}

func (s *Siq) getBackupSiq(tName string) *Siq {
	return s.so.GetBackupSiqByTopic(tName)
}

func (s *Siq) updateBackup(data []byte) error {
	s.backupData = append(s.backupData, data...)
	fmt.Printf("receive backup data %s\n", s.backupData)

	return nil
}

func (s *Siq) consumeAllAt(tName string) {
	t, ok := s.topics[tName]
	if !ok {
		return
	}

	if s.wo.Count() == 0 {
		return
	}

	currentNumber := 0

	for true {
		msg := surface.DequeTopic(t)
		if msg == "" {
			break
		}

		currentNumber, w := s.wo.LockIdleWorker(currentNumber)
		if w == nil {
			fmt.Println("***Idle Not Found****")
			surface.AddTopic(t, msg)
			time.Sleep(1 * time.Second)
			continue
		}
		go func() {
			defer s.wo.Unlock(w)
			do, err := (*w).SendMessage(tName, msg)
			if err != nil {
				fmt.Printf("Worker[%d] fail: %s\n", (*w).GetId(), err)
				s.wo.NotifyFail(w)
			}

			if !do {
				time.Sleep(1 * time.Second)
				do, err = (*w).SendMessage(tName, msg)
				//TODO: pushFirst
				if !do {
					surface.AddTopic(t, msg)
					fmt.Printf("do nothing: %s\n", msg)
				}
			}
			fmt.Printf("finish [%s]\n", msg)

		}()

		currentNumber += 1
	}
}
