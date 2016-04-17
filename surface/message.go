package surface

import "fmt"

// Message is ss
type Message struct {
	id   int
	Text string
}

func (m *Message) Show() {
	fmt.Printf("Message: %s\n", m.Text)
}
