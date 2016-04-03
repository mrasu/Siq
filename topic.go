package siq

import (
	"fmt"
	"mine/Siq/topic_command"
)

// Topic is test
type Topic struct {
	//Name for topic
	Name string
	id int
	messages []Message
}

func NewTopic(n string) (inCh chan topic_command.TopicCommand) {
	inCh = make(chan topic_command.TopicCommand)

	go func() {
		t := Topic{Name: n}
		for {
			command := <- inCh
			fmt.Println(command.Message)

			switch command.Type {
			case topic_command.Add:
				t.AddMessage(command.Message)
				command.ResultChannel <- ""
			case topic_command.Deque:
				command.ResultChannel <- t.Dequeue()
			case topic_command.Show:
				t.Show()
				command.ResultChannel <- ""
			}
		}
	}()

	return inCh
}

func(t *Topic) AddMessage(m string) {
	nm := Message{Text: m}
	t.messages = append(t.messages, nm)
}

func(t *Topic) Dequeue() string {
	if len(t.messages) == 0 {
		return ""
	}

	m := t.messages[0]
	t.messages = t.messages[1:]

	return m.Text
}

func (t *Topic) Show() {
	fmt.Printf("Topic: %s\n", t.Name)
	for _, m := range(t.messages) {
		m.Show()
	}
}