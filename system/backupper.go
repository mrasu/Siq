package system

import (
	"time"
	"encoding/json"
	"github.com/mrasu/Siq/surface"
)

type Backupper struct {
	isPolling          bool
	getTopicsFunc      func() []string
	getMessageFromFunc func(string, int)[] surface.Message
	backupFunc         func(string, []byte) error
	notifyDeadFunc     func(string)

	gatheredMap        map[string]int
	errorCounts        map[string]int
}

func NewBackupper(getTopicFunc func() []string, getMessageFunc func(string, int)[]surface.Message, backupFunc func(string, []byte)error, notifyDeadFunc func(string)) *Backupper {
	b := &Backupper{
		isPolling: false,
		getTopicsFunc: getTopicFunc,
		getMessageFromFunc: getMessageFunc,
		backupFunc: backupFunc,
		notifyDeadFunc: notifyDeadFunc,

		gatheredMap: map[string]int{},
		errorCounts: map[string]int{},
	}

	return b
}

func(b *Backupper) StartPolling() {
	if b.isPolling {
		return
	}

	b.isPolling = true
	go func() {
		for true {
			for _, topic := range(b.getTopicsFunc()) {
				b.backupTopic(topic)
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

func (b *Backupper) backupTopic(topic string) {
	const failThreshold = 1

	fromId, ok := b.gatheredMap[topic]
	if ok == false {
		fromId = -1
	}
	messages := b.getMessageFromFunc(topic, fromId)

	if len(messages) == 0 {
		return
	}
	sendData := b.pack(messages)

	err := b.backupFunc(topic, sendData)
	if err != nil {
		if _, ok := b.errorCounts[topic]; ok == false {
			b.errorCounts[topic] = 1
		} else {
			b.errorCounts[topic] += 1
		}

		if b.errorCounts[topic] > failThreshold {
			b.notifyDeadFunc(topic)
		}
	} else {
		lastMessage := messages[len(messages)-1]
		b.gatheredMap[topic] = lastMessage.Id
		b.errorCounts[topic] = 0
	}
}

func(b *Backupper) pack(messages []surface.Message) []byte {
	json, err := json.Marshal(messages)

	if err != nil {
		panic(err)
	}
	return json
}