package main

import (
	"fmt"
	"github.com/mrasu/Siq"
	"time"
	"github.com/mrasu/Siq/workers"
)

func main() {
	fmt.Println("aa")
	s := siq.NewSiq()
	c := register(s, 1)
	c2 := register(s, 2)
	s.Start()
	//ss := siq.SiqServer{s, 5555}
	//ss.Start()

	c.Callback = func(t string, m string) string {
		fmt.Printf("Cui Worker[%d] error: %s : %s\n", c.GetId(), t, m)
		return "worker is closed"
	}
	c.IsDead = true

	s.Publish("tttt", "mmmm")
	s.Publish("tttt", "mmmm2")
	s.Publish("tttt", "mmmm3")
	s.Publish("tttt", "4444")
	s.Publish("tttt", "+++++5")
	s.Publish("tttt", "^^^^6")
	s.Publish("tttt", "-----7")
	s.Publish("tttt", ":::::8")
	s.Publish("tttt", "======9")

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
			fmt.Printf("Cui Worker[%d] sendmessage: %s : %s\n", id, t, m)
			return ""
		},
		Capacity: 2,
		Id: id,
	}

	s.Register(c)

	return c
}
