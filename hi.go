package hi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	DefaultScheme = "hi"
)

var (
	cli *clientv3.Client
)

type hi struct {
	Endpoints []string
	Scheme    string
}

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

	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   h.Endpoints,
			DialTimeout: 10 * time.Second,
		})
		if err != nil {
			return err
		}
	}

	leaseResp, err := cli.Grant(context.Background(), 12)
	if err != nil {
		log.Fatal(err)
	}

	_, err = cli.Put(context.Background(), fmt.Sprintf("/%s/%s/%s", h.Scheme, name, addr), addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second * 8)
	go func() {
		for t := range ticker.C {
			cli.KeepAliveOnce(context.Background(), leaseResp.ID)
			fmt.Printf("Renew /%s/%s/%s at %s \n", h.Scheme, name, addr, t.Format("15:04:05 2006-01-02"))
		}
	}()

	return nil
}