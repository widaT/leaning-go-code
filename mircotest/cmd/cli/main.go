package main

import (
	"context"
	"fmt"
	hello "mircotest/proto"

	_ "github.com/asim/go-micro/plugins/registry/etcd/v3"
	roundrobin "github.com/asim/go-micro/plugins/wrapper/select/roundrobin/v3"
	"github.com/asim/go-micro/v3"
)

func main() {
	wrapper := roundrobin.NewClientWrapper()
	service := micro.NewService(
		micro.WrapClient(wrapper),
	)
	service.Init()
	cl := hello.NewSayService("wida.micro.srv.greeter", service.Client())
	rsp, err := cl.Hello(context.Background(), &hello.Request{
		Name: "John",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", rsp.Msg)
}
