package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hi"
	"github.com/hi/example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

var (
	etcdLoc = "localhost:2379"
)

func main() {
	hiBuilder := hi.NewResolverBuilder([]string{"localhost:2379"})
	resolver.Register(&hiBuilder)
	rr := balancer.Get("round_robin")
	// grpc.RoundRobin(r)
	var dialAddr = fmt.Sprintf("%s://kai/%s", hiBuilder.Scheme(), "serverA")

	println(dialAddr)

	conn, err := grpc.Dial(dialAddr, grpc.WithBalancerBuilder(rr), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewServerAClient(conn)

	ticker := time.NewTicker(time.Second * 3)
	for {

		req := &pb.HiReq{Name: "kai"}
		resp, err := client.Hi(context.Background(), req)
		println("++++++++++++")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v \n", resp)
		<-ticker.C
	}
}
