package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hi"
	"github.com/hi/example/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	srvName = "serverA"
	port    = ":10014"
)

type serverA struct{}

func (s *serverA) Hi(ctx context.Context, req *pb.HiReq) (*pb.HiResp, error) {

	println("Yeah, it is serverA1 ...")

	return &pb.HiResp{Echo: "I hate you , " + req.Name}, nil
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
	reflection.Register(s)

	println("Hello, I am serverA ...")

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
