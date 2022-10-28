package backend

import (
	"context"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

type activeService struct {
	prefix string
	nodeList map[string]string
	mutex sync.Mutex
}

// 获取服务目录下所有key，初始化到服务的可用节点列表中
func ServiceDiscovery(ctx context.Context, etcdClient *clientv3.Client, serviceTarget string) []string {
	service := &activeService {
		prefix: serviceTarget,
		nodeList: make(map[string]string),
	}
	rangeResp, err := etcdClient.Get(ctx, service.prefix, clientv3.WithPrefix())
	if err != nil {
		log.Error("can not get service KV: ", err)
		return nil
	}

	service.mutex.Lock()
	for _, kv := range rangeResp.Kvs {
		service.nodeList[string(kv.Key)] = string(kv.Value)
	}
	service.mutex.Unlock()

	go watchServiceUpdate(etcdClient, service)

	var result []string
	for _, val := range service.nodeList {
		result = append(result, val)
	}
	return result
}

// 监控服务目录下的事件
func watchServiceUpdate(etcdClient *clientv3.Client, service *activeService) {
	watcher := clientv3.NewWatcher(etcdClient)
	// Watch 服务目录下的更新
	watchChan := watcher.Watch(context.TODO(), service.prefix, clientv3.WithPrefix())
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			service.mutex.Lock()
			switch int(event.Type) {
			case 0://PUT事件，目录下有了新key
				service.nodeList[string(event.Kv.Key)] = string(event.Kv.Value)
			case 1://DELETE事件，目录中有key被删掉(Lease过期，key 也会被删掉)
				delete(service.nodeList, string(event.Kv.Key))
			}
			service.mutex.Unlock()
		}
	}
}
