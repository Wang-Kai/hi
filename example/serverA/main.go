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

const (
	srvName = "serverA"
	port    = ":10013"
)

type serverA struct{}

func (s *serverA) Hi(ctx context.Context, req *pb.HiReq) (*pb.HiResp, error) {
	println("Yeah, it is serverA ...")

	return &pb.HiResp{Echo: "Hi " + req.Name + ", this response comes from ServerA"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	// register serverA to etcd
	h := hi.NewHi([]string{"localhost:2379"}, "hi")
	err = h.Register(srvName, fmt.Sprintf("127.0.0.1%s", port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterServerAServer(s, &serverA{})

	println("Hello, I am serverA ...")

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
