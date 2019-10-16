package main

import (
	"io"
	"log"
	"net"

	"context"
	"github.com/widaT/gorpc_demo/grpc/pb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchReply, error) {
	return &pb.SearchReply{Name:"hello " + in.GetName()}, nil
}

func (s *server) Search2(stream pb.Searcher_Search2Server) error {

	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &pb.SearchReply{ Name:"hello:" + args.GetName()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}

	return nil
}



func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSearcherServer(s, &server{})
	s.Serve(lis)
}