package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type Param struct {
	A,B int
}

func main() {
	conn, err := net.Dial("tcp", "localhost:9700")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	client :=rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	//同步调用
	var reply int
	err = client.Call("test.Multiply", Param{34,35}, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
	//异步调用
	done := make(chan *rpc.Call, 1)
	client.Go("test.Divide", Param{34,17}, &reply,done)
	select {
		case d := <-done:
			fmt.Println(* d.Reply.(*int))

		case <-time.After(3e9):
			fmt.Println("time out")
	}
}