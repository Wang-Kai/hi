package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Wang-Kai/hi"
	"github.com/Wang-Kai/hi/example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	etcdLoc = "localhost:2379"
)
var ConnMap = make(map[string]*grpc.ClientConn, 2)

func init() {
	// register resolver
	hiBuilder := hi.NewResolverBuilder([]string{"localhost:2379"})
	resolver.Register(&hiBuilder)

}

func main() {
	// build connection of serverB
	serverBConn, err := grpc.Dial("hi://author/serverB", grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	if err != nil {
		log.Fatal(err)
	}
	println(serverBConn)
	ConnMap["serverB"] = serverBConn

	// build connection of serverA
	serverAConn, err := grpc.Dial("hi://author/serverA", grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	if err != nil {
		log.Fatal(err)
	}
	println(serverAConn)
	ConnMap["serverA"] = serverAConn
	defer func() {
		for _, cc := range ConnMap {
			cc.Close()
		}
	}()

	for range time.Tick(time.Second * 3) {
		println("Call serverB")
		serverBClient := pb.NewServerBClient(ConnMap["serverB"])

		println("New Client")
		helloResp, err := serverBClient.Hello(context.Background(), &pb.HelloReq{Name: "China"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v \n", helloResp)

		println("Call serverA")
		serverAClient := pb.NewServerAClient(ConnMap["serverA"])
		hiResp, err := serverAClient.Hi(context.Background(), &pb.HiReq{Name: "kai"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v \n", hiResp)
	}
}
