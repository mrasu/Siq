package siq

import (
	"fmt"
	"mine/Siq/topic_command"
)

// Siq is distributed job queue
type Siq struct {
	topics map[string]chan topic_command.TopicCommand
}

func NewSiq() Siq {
	return Siq{topics: map[string] chan topic_command.TopicCommand{}}
}

// Add is text
func(s *Siq) Publish(tName string, m string) {
	s.addMessage(tName, m)
}

func(s *Siq) addMessage(tName string, m string) {
	topic := s.getTopicChannel(tName)
	<- topic_command.AddTopic(topic, m)
}

func(s *Siq) Consume(tName string) string {
	var t chan topic_command.TopicCommand
	t, ok := s.topics[tName]
	if !ok {
		return ""
	}

	return <- topic_command.DequeTopic(t)
}

func(s *Siq) getTopicChannel(n string) chan topic_command.TopicCommand {
	if c, ok := s.topics[n]; ok {
		return c
	}
	c := NewTopic(n)
	s.topics[n] = c
	return c
}

func(s *Siq) Show() {
	fmt.Println("Show Topic!!")
	for _, t := range(s.topics) {
		<- topic_command.ShowTopic(t)
	}
}