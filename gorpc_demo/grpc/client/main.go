package main

import (
	"context"
	"fmt"
	"github.com/widaT/gorpc_demo/grpc/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSearcherClient(conn)

	stream, err := c.Search2(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			if err := stream.Send(&pb.SearchRequest{Name: "world"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()

	for {
		rep, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(rep.GetName())
	}

}
