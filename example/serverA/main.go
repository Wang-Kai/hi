package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Wang-Kai/hi"
	"github.com/Wang-Kai/hi/example/pb"

	"google.golang.org/grpc"
)

var (
	svcName *string
	port    *string
)

func init() {
	svcName = flag.String("name", "serverA", "The name of microservice")
	port = flag.String("port", "10013", "The port of microservice")
	flag.Parse()
}

type serverA struct{}

func (s *serverA) Hi(ctx context.Context, req *pb.HiReq) (*pb.HiResp, error) {
	println("Yeah, it is serverA ...")

	return &pb.HiResp{Echo: "Hi " + req.Name + ", this response comes from ServerA"}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatal(err)
	}

	// register serverA to etcd
	h := hi.NewHi([]string{"localhost:2379"}, "hi")
	err = h.Register(*svcName, fmt.Sprintf("127.0.0.1:%s", *port))
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
