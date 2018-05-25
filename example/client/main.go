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

func main() {
	hiBuilder := hi.NewResolverBuilder([]string{"localhost:2379"})
	resolver.Register(&hiBuilder)
	var dialAddr = fmt.Sprintf("%s://foo/%s", hiBuilder.Scheme(), "serverA")

	conn, err := grpc.Dial(dialAddr, grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewServerAClient(conn)

	for range time.Tick(time.Second * 3) {
		req := &pb.HiReq{Name: "Foo"}

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelFunc()

		resp, err := client.Hi(ctx, req)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v \n", resp)
	}
}
