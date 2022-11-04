package utils

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"time"
)

var (
	AllKeysCounts = "AllKeysCounts"
)

func PutData(cli *clientv3.Client, key string, value string, requestTimeout time.Duration) (putResponse *clientv3.PutResponse) {

	// 创建一个90天的租约，租约的第二个参数是以秒为单位的

	resp, err := cli.Grant(context.TODO(), int64(60*60*24*90))
	if err != nil {
		log.Fatal(err)
	}
	// 5秒钟之后, /lmh/ 这个key就会被移除
	putResponse, err = cli.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	return
}

func GetData(cli *clientv3.Client, key string, requestTimeout time.Duration) (getResponse *clientv3.GetResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		fmt.Println("err = ", err)
	}
	cancel()
	return resp
}

//根据key前缀获取数据
func GetDataPrefix(cli *clientv3.Client, key string, requestTimeout time.Duration) (getResponse *clientv3.GetResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("err = ", err)
	}
	cancel()
	return resp
}

//统计所有key数量
func AllKeysCount(cli *clientv3.Client, key string, num int, requestTimeout time.Duration) {

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, key)
	if err != nil {
		fmt.Println("err = ", err)
	}
	cancel()

	allKeysCountStr := "0"
	for _, ev := range resp.Kvs {
		allKeysCountStr = string(ev.Value)
	}
	allKeysCountInt, _ := strconv.Atoi(allKeysCountStr)

	_, err = cli.Put(context.TODO(), key, strconv.Itoa(allKeysCountInt+num))
	if err != nil {
		log.Fatal(err)
	}
}
