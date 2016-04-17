package surface

import (
	"fmt"
)

// Topic is test
type Topic struct {
	//Name for topic
	Name     string
	id       int
	messages []Message
}

func NewTopic(n string) chan *TopicCommand {
	inCh := make(chan *TopicCommand)

	go func() {
		t := Topic{Name: n}
		for {
			command := <-inCh

			switch command.Type {
			case Add:
				t.AddMessage(command.Message)
				command.ResultChannel <- ""
			case Deque:
				command.ResultChannel <- t.Dequeue()
			case Show:
				t.Show()
				command.ResultChannel <- ""
			}
		}
	}()

	return inCh
}

func (t *Topic) AddMessage(m string) {
	nm := Message{Text: m}
	t.messages = append(t.messages, nm)
}

func (t *Topic) Dequeue() string {
	if len(t.messages) == 0 {
		return ""
	}

	m := t.messages[0]
	t.messages = t.messages[1:]

	return m.Text
}

func (t *Topic) Show() {
	fmt.Printf("Topic: %s\n", t.Name)
	for _, m := range t.messages {
		m.Show()
	}
}
