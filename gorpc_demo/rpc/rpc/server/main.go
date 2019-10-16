package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
)

type Param struct {
	A,B int
}
type Server struct {}

func (t *Server) Multiply(args *Param, reply *int) error {
	*reply = args.A * args.B
	return nil
}
func (t *Server) Divide(args *Param, reply *int) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	*reply = args.A / args.B
	return nil
}

func main() {
	rpc.RegisterName("test", new(Server))
	listener, err := net.Listen("tcp", ":9700")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	for  {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeConn(conn)
	}
}
