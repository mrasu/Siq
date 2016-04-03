package main

import (
	"fmt"
	"mine/Siq"
)

func main() {
	fmt.Println("aa")
	s := siq.NewSiq()
	s.Publish("tttt", "mmmm")
	s.Publish("tttt", "mmmm2")
	s.Publish("tttt", "mmmm3")
	s.Publish("tttt", "4444")
	s.Publish("tttt", "+++++")

	s.Show()

	for {
		a := s.Consume("tttt")
		if a == "" {
			break
		}
		fmt.Println(a)

	}
}
