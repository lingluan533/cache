package backend
/*
import (
	"context"
	"fmt"
clientv3 "go.etcd.io/etcd/client/v3"
	"github.com/labstack/gommon/log"

	"time"
)

type Service struct {
	etcdClient *clientv3.Client
	key string
	val string
	leaseId clientv3.LeaseID
}

// 指定client端，Endpoints是etcd server的机器列表，DialTimeout是计算节点链接服务的超时时间
func NewService(endpoints []string, key string, val string, lease int64) (*Service, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Error("etcdClient start error: ", err)
		return nil, err
	}

	service := &Service{
		etcdClient: etcdClient,
		key: key,
		val: val,
	}



	return service, nil
}

//注册服务
func (this *Service) RegService(id string, name string, address string) error {
	kv := clientv3.NewKV(this.client)
	key_prefix := "/etcd_services/"
	ctx := context.Background()
	lease := clientv3.NewLease(this.client)

	//设置租约过期时间为20秒
	leaseRes, err := clientv3.NewLease(this.client).Grant(ctx, 20)
	if err != nil {
		return err
	}
	_, err = kv.Put(context.Background(), key_prefix+id+"/"+name, address, clientv3.WithLease(leaseRes.ID)) //把服务的key绑定到租约下面
	if err != nil {
		return err
	}
	//续租时间大概自动为租约的三分之一时间，context.TODO官方定义为是你不知道要传什么
	keepaliveRes, err := lease.KeepAlive(context.TODO(), leaseRes.ID)
	if err != nil {
		return err
	}
	go lisKeepAlive(keepaliveRes)
	return err
}

func lisKeepAlive(keepaliveRes <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case ret := <-keepaliveRes:
			if ret != nil {
				fmt.Println("续租成功", time.Now())
			}
		}
	}
}
*/