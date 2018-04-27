package hi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/resolver"
)

// Implement resolver.Resolver
type Resolver struct {
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

	b.cc = cc

	r := &Resolver{}

	log.Infof("Watch ==> %s \n", fmt.Sprintf("/%s/%s/\n", target.Scheme, target.Endpoint))
	go b.watch(fmt.Sprintf("%s/%s/", target.Scheme, target.Endpoint))

	return r, nil
}

func (b *Builder) watch(keyPrefix string) {
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

	b.cc.NewAddress(addrList)

	// start to watch keys which prefix with keyPrefix
	wch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for wresp := range wch {
		for _, ev := range wresp.Events {
			evKey := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)

			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, evKey) {
					addrList = append(addrList, resolver.Address{Addr: evKey})
					b.cc.NewAddress(addrList)
				}
			case mvccpb.DELETE:
				if list, ok := remove(addrList, evKey); ok {
					addrList = list
					b.cc.NewAddress(addrList)
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
