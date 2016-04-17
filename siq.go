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
	wm             *WorkerObserver
}

func NewSiq() *Siq {
	s := &Siq{
		topics: map[string]chan *surface.TopicCommand{},
		wm:     newWorkerManager(),
	}

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
	s.wm.Add(&w)
}

func (s *Siq) Unregister(w workers.Worker) {
	s.wm.Delete(&w)
}

func (s *Siq) Publish(tName string, m string) {
	s.addMessage(tName, m)
	s.StartConsume(tName)
}

func (s *Siq) addMessage(tName string, m string) {
	topic := s.getTopicChannel(tName)
	surface.AddTopic(topic, m)
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
	c := surface.NewTopic(n)
	s.topics[n] = c
	return c
}

func (s *Siq) Show() {
	fmt.Println("Show Topic!!")
	for _, t := range s.topics {
		surface.ShowTopic(t)
	}
}

func (s *Siq) consumeAllAt(tName string) {
	t, ok := s.topics[tName]
	if !ok {
		return
	}

	if s.wm.Count() == 0 {
		return
	}

	currentNumber := 0

	for true {
		msg := surface.DequeTopic(t)
		if msg == "" {
			break
		}

		currentNumber, w := s.wm.LockIdleWorker(currentNumber)
		if w == nil {
			fmt.Println("***Idle Not Found****")
			surface.AddTopic(t, msg)
			time.Sleep(1 * time.Second)
			continue
		}
		go func() {
			defer s.wm.Unlock(w)
			do, err := (*w).SendMessage(tName, msg)
			if err != nil {
				fmt.Printf("Worker[%d] fail: %s\n", (*w).GetId(), err)
				s.wm.NotifyFail(w)
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
