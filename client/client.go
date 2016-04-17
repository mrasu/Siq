package client

import (
	"errors"
	"fmt"
	"github.com/nobonobo/jsonrpc"
	"net/http"
)

type SiqClient struct {
	Callback func(ClientChannel) string
}

type ClientChannel struct {
	Name    string
	Message string
}

func (c *SiqClient) Go(args *[]string, reply *int) error {
	fmt.Println("Client Channel Go!!")
	fmt.Println(*args)

	a := *args

	data := ClientChannel{
		Name:    a[0],
		Message: a[1],
	}

	result := c.Callback(data)

	return errors.New(result)
}

func (s *SiqClient) Start() {
	server := jsonrpc.NewServer()
	server.Register(s)
	http.Handle(jsonrpc.DefaultRPCPath, server)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
