package workers

import (
	"errors"
	"fmt"
)

type CuiWorker struct {
	Id int
	Callback         func(string, string) string
	Capacity         int
	IsDead bool
	currentCallCount int
}

func (c *CuiWorker) SendMessage(topicName string, messageText string) (bool, error) {
	if c.isSendable() == false {
		return false, nil
	}

	c.currentCallCount += 1
	em := c.Callback(topicName, messageText)
	c.currentCallCount -= 1
	if em != "" {
		return true, errors.New(em)
	}

	return true, nil
}

func (c *CuiWorker) isSendable() bool {
	return c.currentCallCount < c.Capacity
}

func (c *CuiWorker) Ping() bool {
	fmt.Println("Ping comming")
	return !c.IsDead
}

func (c *CuiWorker) GetId() int {
	return c.Id
}