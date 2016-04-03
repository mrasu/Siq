package topic_command

const (
	Add TopicCommandType = iota
	Show
	Deque
)

type TopicCommandType int

type TopicCommand struct {
	Type          TopicCommandType
	Message       string
	// To Result Command
	ResultChannel chan string
}

func AddTopic(ch chan TopicCommand, message string) chan string{
	c := TopicCommand{Type: Add, Message: message}
	return c.send(ch)
}

func ShowTopic(ch chan TopicCommand) chan string{
	c := TopicCommand{Type: Show}
	return c.send(ch)
}

func DequeTopic(ch chan TopicCommand) chan string{
	c := TopicCommand{Type: Deque}
	return c.send(ch)
}

func(c *TopicCommand) send(ch chan TopicCommand) chan string {
	resCh := make(chan string)
	c.ResultChannel = resCh

	ch <- *c
	return resCh
}