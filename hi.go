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
	"strings"
	"time"

	"google.golang.org/grpc/resolver"

	"github.com/coreos/etcd/clientv3"
)

const (
	DefaultScheme = "hi"
)

var (
	cli *clientv3.Client
)

func parseTarget(target string) (ret resolver.Target) {
	var ok bool
	ret.Scheme, ret.Endpoint, ok = split2(target, "://")
	if !ok {
		return resolver.Target{Endpoint: target}
	}
	ret.Authority, ret.Endpoint, _ = split2(ret.Endpoint, "/")
	return ret
}

func split2(s, sep string) (string, string, bool) {
	spl := strings.SplitN(s, sep, 2)
	if len(spl) < 2 {
		return "", "", false
	}
	return spl[0], spl[1], true
}

type hi struct {
	Endpoints []string
	Scheme    string
}

// NewHi create a Hi instance
func NewHi(endpoints []string, scheme string) hi {
	if scheme == "" {
		scheme = DefaultScheme
	}

	return hi{Endpoints: endpoints, Scheme: scheme}
}

// Unregiste delete name from etcd
func (h *hi) Unregister(name string) error {
	_, err := cli.Delete(context.Background(), name, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}

// Register register microserver address under the scheme/name
func (h *hi) Register(name, addr string) error {
	var err error

	// create client if not define
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   h.Endpoints,
			DialTimeout: 10 * time.Second,
		})
		if err != nil {
			return err
		}
	}

	// create lease with 12s TTL
	leaseResp, err := cli.Grant(context.Background(), 12)
	if err != nil {
		return err
	}

	// put key & value
	key := fmt.Sprintf("%s/%s/%s", h.Scheme, name, addr)
	_, err = cli.Put(context.Background(), key, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	// generate a goroutine to keep lease alive per 10s
	go func() {
		for range time.Tick(time.Second * 10) {
			cli.KeepAliveOnce(context.Background(), leaseResp.ID)
		}
	}()

	return nil
}
