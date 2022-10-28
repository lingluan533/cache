package backend

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
)

type Resolver struct {
	sync.RWMutex
	Client    *clientv3.Client
	cc        resolver.ClientConn
	prefix    string
	addresses map[string]resolver.Address
}

func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {
	// todo
}

func (r *Resolver) Close() {
	// todo
}

func (r *Resolver) watcher() {
	r.addresses = make(map[string]resolver.Address)

	response, err := r.Client.Get(context.Background(), r.prefix, clientv3.WithPrefix())

	if err == nil {
		for _, kv := range response.Kvs {
			r.setAddress(string(kv.Key), string(kv.Value))
		}

		r.cc.UpdateState(resolver.State{
			Addresses: r.getAddresses(),
		})
	}

	watch := r.Client.Watch(context.Background(), r.prefix, clientv3.WithPrefix())

	for response := range watch {
		for _, event := range response.Events {
			switch event.Type {
			case mvccpb.PUT:
				r.setAddress(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				r.delAddress(string(event.Kv.Key))
			}
		}

		r.cc.UpdateState(resolver.State{
			Addresses: r.getAddresses(),
		})
	}
}

func (r *Resolver) setAddress(key, address string) {
	r.Lock()
	defer r.Unlock()
	r.addresses[key] = resolver.Address{Addr: string(address)}
}

func (r *Resolver) delAddress(key string) {
	r.Lock()
	defer r.Unlock()
	delete(r.addresses, key)
}

func (r *Resolver) getAddresses() []resolver.Address {
	addresses := make([]resolver.Address, 0, len(r.addresses))

	for _, address := range r.addresses {
		addresses = append(addresses, address)
	}

	return addresses
}
