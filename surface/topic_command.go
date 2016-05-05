package surface

import (
	"math/rand"
	"encoding/json"
	"strconv"
)

const (
	Add TopicCommandType = iota
	Show
	Deque
	GetFrom
)

type TopicCommandType int

type TopicCommand struct {
	Type    TopicCommandType
	Message string
	// To Result Command
	ResultChannel chan string
	id            int
}

func AddTopic(ch chan *TopicCommand, message string) string {
	c := newTopicCommand(Add)
	c.Message = message
	return <-c.send(ch)
}

func ShowTopic(ch chan *TopicCommand) string {
	c := newTopicCommand(Show)
	return <-c.send(ch)
}

func DequeTopic(ch chan *TopicCommand) string {
	c := newTopicCommand(Deque)
	return <-c.send(ch)
}

func GetTopicFrom(ch chan *TopicCommand, fromId int) []Message {
	c := newTopicCommand(GetFrom)
	c.Message = strconv.Itoa(fromId)
	messageJson := <- c.send(ch)

	var messages []Message
	err := json.Unmarshal([]byte(messageJson), &messages)
	if err != nil {
		panic(err)
	}
	return messages
}

func newTopicCommand(t TopicCommandType) *TopicCommand {
	c := &TopicCommand{Type: t}
	c.id = rand.Int()
	return c
}

func (c *TopicCommand) send(ch chan *TopicCommand) chan string {
	resCh := make(chan string)
	c.ResultChannel = resCh

	ch <- c

	return resCh
}
