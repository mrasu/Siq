package main

import (
	"fmt"
	"github.com/mrasu/Siq"
	"time"
	"github.com/mrasu/Siq/workers"
)

func main() {
	sList := siq.StartSiq()
	s1, s2 := sList[0], sList[1]
	c := register(s1, 1)
	c2 := register(s2, 2)
	//ss := siq.SiqServer{s, 5555}
	//ss.Start()

	c.Callback = func(t string, m string) string {
		fmt.Printf("Cui Worker[%d] error: %s : %s\n", c.GetId(), t, m)
		return "worker is closed"
	}
	c.IsDead = true


	currentSiq := s1
	for i := 0; i < 10; i++ {
		otherSiq := currentSiq.Publish("tttt", fmt.Sprint("ttt%d", i))
		if otherSiq != nil {
			currentSiq = otherSiq
		}
	}

	currentSiq = s2
	for i := 0; i < 10; i++ {
		otherSiq := currentSiq.Publish("uuuu", fmt.Sprintf("uuuu%d", i))
		if otherSiq != nil {
			currentSiq = otherSiq
		}
	}

	go func() {
		//time.Sleep(3 * time.Second)
		c.Callback = func(t string, m string) string {
			fmt.Printf("Cui Worker[%d] error: %s : %s\n", c.GetId(), t, m)
			return "worker is closed"
		}
		c.IsDead = true
	}()
	go func() {
		time.Sleep(3 * time.Second)
		c2.Callback = func(t string, m string) string {
			fmt.Printf("Cui Worker[%d] error: %s : %s\n", c2.GetId(), t, m)
			return "worker is closed"
		}
		c2.IsDead = true
	}()
	ch := make(chan int)
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 1
	}()
	<-ch
}

func register(s *siq.Siq, id int) *workers.CuiWorker {
	c := &workers.CuiWorker{
		Callback: func(t string, m string) string {
			time.Sleep(1 * time.Second)
			fmt.Printf("Cui Worker[%d] receiveMessage: %s : %s\n", id, t, m)
			return ""
		},
		Capacity: 2,
		Id: id,
	}

	s.Register(c)

	return c
}
