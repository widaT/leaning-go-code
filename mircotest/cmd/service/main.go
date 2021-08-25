package main

import (
	"context"
	"log"
	hello "mircotest/proto"
	"os"

	_ "github.com/asim/go-micro/plugins/registry/etcd/v3" //v3版本registry变成 plugins 需要import
	"github.com/asim/go-micro/v3"
)

type Hello struct{}

func (s *Hello) Hello(ctx context.Context, req *hello.Request, rsp *hello.SayResponse) error {
	log.Print("Received Say.Hello request")
	hostname, _ := os.Hostname()
	rsp.Msg = "Hello " + req.Name + " ,Im " + hostname
	rsp.Header = make(map[string]*hello.Pair)
	rsp.Header["name"] = &hello.Pair{Key: 1, Values: "abc"}
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("wida.micro.srv.greeter"),
	)
	service.Init()

	// Register Handlers
	hello.RegisterSayHandler(service.Server(), new(Hello))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
