/*
 * Created on Fri May 25 2018
 *
 * Copyright (c) WangKai
 * Email persistence201306@gmail.com
 */

package hi

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
)

// Implement resolver.Resolver
type Resolver struct {
	cc resolver.ClientConn
}

func (r *Resolver) Close() {
	println("Call Close method")
}

func (r *Resolver) ResolveNow(opt resolver.ResolveNowOption) {
	println("Call ResolveNow method")

}

func NewResolverBuilder(etcdEndPoints []string) Builder {
	return Builder{
		endPoints: etcdEndPoints,
	}
}

// Implement resolver.Builder
type Builder struct {
	cc        resolver.ClientConn
	endPoints []string
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var err error

	// 建立对某一个 target 的解析
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   b.endPoints,
			DialTimeout: 10 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}

	r := &Resolver{
		cc: cc,
	}

	log.Printf("Watch ==> %s \n", fmt.Sprintf("%s/%s/\n", target.Scheme, target.Endpoint))
	go r.watcher(fmt.Sprintf("%s/%s/", target.Scheme, target.Endpoint))

	return r, nil
}

func (r *Resolver) watcher(keyPrefix string) {
	var addrList []resolver.Address

	// first get all address under this keyPrefix
	resp, err := cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
	}

	for _, kv := range resp.Kvs {
		addr := resolver.Address{Addr: strings.TrimPrefix(string(kv.Key), keyPrefix)}
		addrList = append(addrList, addr)
	}

	fmt.Printf("%+v \n", addrList)
	r.cc.NewAddress(addrList)

	// start to watch keys which prefix with keyPrefix
	wch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for wresp := range wch {
		for _, ev := range wresp.Events {
			evKey := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)

			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, evKey) {
					addrList = append(addrList, resolver.Address{Addr: evKey})
					r.cc.NewAddress(addrList)
					println("Add an connection: " + evKey)
				}
			case mvccpb.DELETE:
				if list, ok := remove(addrList, evKey); ok {
					addrList = list
					r.cc.NewAddress(addrList)
					println("Remove an connnection: " + evKey)
				}
			}
		}
	}
}

func (b *Builder) Scheme() string {
	return DefaultScheme
}

func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
