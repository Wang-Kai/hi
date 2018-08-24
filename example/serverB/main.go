package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Wang-Kai/hi"
	"github.com/Wang-Kai/hi/example/pb"
	"google.golang.org/grpc"
)

var (
	srvName = "serverB"
	port    = ":10015"
)

type serverB struct{}

func (s *serverB) Hello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	println("Yeah, it is serverB ...")

	return &pb.HelloResp{Echo: "Hello" + req.Name + " ..."}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	h := hi.NewHi([]string{"localhost:2379"}, "hi")
	err = h.Register(srvName, fmt.Sprintf("127.0.0.1%s", port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterServerBServer(s, &serverB{})

	println("I am serverB ...")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
