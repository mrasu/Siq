package siq

import (
	"errors"
	"fmt"
	"github.com/nobonobo/jsonrpc"
	"net/http"
)

type SiqServer struct {
	S    Siq
	Port int
}

func (ss *SiqServer) Publish(args *[]string, reply *int) error {
	a := *args
	fmt.Printf("Receive Publish %s\n", a)
	ss.S.Publish(a[0], a[1])

	return errors.New("")
}

func (ss *SiqServer) Start() {
	server := jsonrpc.NewServer()
	server.Register(ss)
	http.Handle(jsonrpc.DefaultRPCPath, server)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
