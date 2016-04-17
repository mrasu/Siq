package workers

// Worker is Interface for rpc client and cui
type Worker interface {
	GetId() int
	SendMessage(topicName string, messageText string) (bool, error)
	Ping() bool
}
