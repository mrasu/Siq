package surface

import "fmt"

type Message struct {
	Id   int
	Text string
}

func (m *Message) Show() {
	fmt.Printf("Message: %s\n", m.Text)
}
