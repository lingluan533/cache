package backend

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"

	"cache/dataStruct"
	"net"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

func NewRedisBackend(ctx context.Context, config *dataStruct.GlobalConfig) (*redis.Client, error) {
	//go StartGrpcPort(":" + "8880")

	if config.Cache.CommonConfig.SyncInternal == 0 || config.Cache.CommonConfig.SyncSizeLimit == 0 {
		log.Error("internal or limit is zero")
		return nil, fmt.Errorf("internal or limit is zero")
	}
	//根据配置文件的信息来连接rdb
	rdb := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(config.Cache.RedisGroup["redis1"].Host, config.Cache.RedisGroup["redis1"].Port),
		Password: config.Cache.CommonConfig.Password,
		DB:       config.Cache.CommonConfig.DB,
	})
	//给五类账本分别开启订阅
	for _, t := range LedgerMap {
		client := &subscribeClient{
			rdb:         rdb,
			config:      config,
			ledgerName:  t,
			currentSize: 0,
			currentNum:  0,
		}
		//t需要传入闭包或者重新声明
		go func(t string) {
			log.Info("open subscribe client for: ", t)
			if err := initDataTypeSubscribe(ctx, client); err != nil {
				log.Error("initDataTypeSubscribe error: ", err)
				_ = rdb.Close()
				return
			}

		}(t)
	}

	return rdb, nil
}

// var wg sync.WaitGroup
//var wg2 sync.WaitGroup

type subscribeClient struct {
	wg, wg2     sync.WaitGroup
	rdb         *redis.Client
	config      *dataStruct.GlobalConfig
	ledgerName  string //账本类型
	currentSize int
	currentNum  int
}

func (s *subscribeClient) RedisCacheKey() string {
	return "_cache" + s.ledgerName
}

func (s *subscribeClient) FailedDataKey() string {
	return "_failed" + s.ledgerName
}

func initDataTypeSubscribe(ctx context.Context, client *subscribeClient) error {
	//发布订阅，根据不同的账本名.通道的名字是五类账本命名的，所以发布者要用这个名字来发布订阅
	pubsub := client.rdb.Subscribe(ctx, client.ledgerName)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Error("pubsub receive error: ", err)
		return err
	}
	//？？etcd客户端

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{net.JoinHostPort(client.config.Consensus.EtcdGroup[client.config.Common.LedgerName[client.ledgerName].Leader].Host,
			client.config.Consensus.EtcdGroup[client.config.Common.LedgerName[client.ledgerName].Leader].Port)},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Error("etcd client start error: ", err)
	}

	defer etcdClient.Close()

	// Go channel which receives messages.
	ch := pubsub.Channel()
	// chch := pubsub.Channel()
	// for msg := range chch {

	// 	msg := msg
	// 	client.currentSize += len(msg.Payload)

	// 	if err := client.rdb.LPush(ctx, client.RedisCacheKey(), msg.Payload).Err(); err != nil {

	// 		log.Error("LPush error: ", err)
	// 	}
	// 	//设置过期时间，如果发送数据超时，会定时清理数据
	// 	if err != nil {

	// 		log.Error("TTL error: ", err)
	// 	}

	// }

	//协程，上传数据的时间限制，固定时间段内上传一次数据
	// go func() {
	// 	ticker := time.NewTicker(time.Second * time.Duration(client.config.Cache.CommonConfig.SyncInternal))

	// 	defer ticker.Stop()
	// 	log.Info("按照时间推送") //判断有没有数据需要推送
	// 	for {
	// 		<-ticker.C
	// 		wg.Add(1)
	// 		//log.Info("消息总长度 ", client.currentSize)
	// 		wg2.Wait()
	// 		//start := time.Now().UnixNano()
	// 		//a := client.currentSize
	// 		if err := sendDataBlock(ctx, client, etcdClient); err != nil {
	// 			Status.FailPushNumber++
	// 			log.Error("sendDataBlock error: ", err)
	// 			//return
	// 		}
	// 		//end := time.Now().UnixNano()
	// 		//fmt.Printf("endtime:%v\n", end)
	// 		//log.Infof("按时间推送%v大小的数据总用时：%v ms\n", a, (end-start)/1000000)
	// 		//记录按时间限制推送的数据次数
	// 		Status.TimePushNumber++
	// 		wg.Done()
	// 	}
	// }()
	ch2 := make(chan error)
	ch3 := make(chan int)
	defer close(ch2)
	defer close(ch3)
	// go func() error {
	// 	for {
	// 		select {
	// 		case err := <-ch2:
	// 			return err //再穿一个给主函数，让主函数return
	// 		case x := <-ch3:
	// 			//if x > client.config.Cache.CommonConfig.SyncSizeLimit {
	// 			if x > client.config.Cache.CommonConfig.SyncSizeLimit*1024*1024 {
	// 				// 发送
	// 				wg.Add(1)
	// 				log.Info("消息总长度 ", client.currentSize)
	// 				//a := client.currentSize
	// 				log.Info("按照数据大小阈值推送")
	// 				wg2.Wait()
	// 				//start := time.Now().UnixNano()
	// 				if err := sendDataBlock(ctx, client, etcdClient); err != nil {
	// 					Status.FailPushNumber++
	// 					log.Error("sendDataBlock error: ", err)
	// 					return err
	// 				}
	// 				//end := time.Now().UnixNano()
	// 				//fmt.Printf("endtime:%v\n", end)
	// 				//log.Infof("按阈值推送%v大小的数据总用时：%v ms\n", a, (end-start)/1000000)
	// 				Status.SizePushNumber++
	// 				//记录按消息总大小限制推送的数据次数
	// 				//这里可以发一个通道给下面的通道的go func 只有接收到才能执行，但是这样又只能有一个线程
	// 				wg.Done()
	// 			}

	// 		}

	// 	}
	// }()

	go func() error {
		for {
			select {
			case err := <-ch2:
				return err //再穿一个给主函数，让主函数return
			case x := <-ch3:
				//if x > client.config.Cache.CommonConfig.SyncSizeLimit {
				if x >= 20 {
					// 发送

					log.Info("消息总长度 ", client.currentSize)
					//a := client.currentSize
					log.Info("按照数据大小阈值推送类型", client.ledgerName)
					client.wg2.Wait()
					//start := time.Now().UnixNano()
					if err := sendDataBlock(ctx, client, etcdClient); err != nil {
						Status.FailPushNumber++
						log.Error("sendDataBlock error: ", err)
						return err
					}
					//end := time.Now().UnixNano()
					//fmt.Printf("endtime:%v\n", end)
					//log.Infof("按阈值推送%v大小的数据总用时：%v ms\n", a, (end-start)/1000000)
					Status.SizePushNumber++
					//记录按消息总大小限制推送的数据次数
					//这里可以发一个通道给下面的通道的go func 只有接收到才能执行，但是这样又只能有一个线程
					client.wg.Done()
				}

			}

		}
	}()
	// Consume messages.获取订阅的消息payload是消息，channel是通道名
	for msg := range ch {
		client.wg.Wait()
		//计算client获取到消息的总长度
		msg := msg
		client.currentSize += len(msg.Payload)
		client.currentNum++
		//ch3 <- client.currentSize
		if client.currentNum == 20 {
			client.wg.Add(1)
		}

		ch3 <- client.currentNum

		//log.Info("消息总长度 ", client.currentSize)
		go func() {
			// if client.currentSize > client.config.Cache.CommonConfig.SyncSizeLimit*1024*1024 {
			// 	client.m.Lock()
			// 	defer client.m.Unlock()
			// }
			client.wg2.Add(1)
			//可以使用LPush()方法将数据从左侧压入链表，client.RedisCacheKey()是key,payload是value

			if err := client.rdb.LPush(ctx, client.RedisCacheKey(), msg.Payload).Err(); err != nil {

				log.Error("LPush error: ", err)
				ch2 <- err
			}
			//设置过期时间，如果发送数据超时，会定时清理数据

			res, err := client.rdb.Expire(ctx, client.RedisCacheKey(), time.Minute*60).Result()
			client.wg2.Done()
			if err != nil {

				log.Error("TTL error: ", err)
				ch2 <- err
			}
			if res {
				//	log.Info("设置TTL成功")
			} else {
				log.Info("设置TTL失败")
			}

		}()

	}

	return nil
}

func sendDataBlock(ctx context.Context, client *subscribeClient, etcdClient *clientv3.Client) error {
	// client.m.Lock()

	// todo: 应该分次小量取出。LRange():获取某个选定范围的元素集 0 -1表示全部元素
	results, err := client.rdb.LRange(ctx, client.RedisCacheKey(), 0, -1).Result()
	if err != nil {
		log.Error("LRange error: ", err)
		return err
	}

	// 数据推送
	if len(results) != 0 {
		//log.Info("legderName: ", client.ledgerName, " 待发送数据: ", results)
		log.Info("start to process data...")
		if err := NewRpcClient(ctx, client.config, client.ledgerName, etcdClient, results); err != nil {
			log.Error("rpcClient error: ", err)
			return err
		}
	}
	//推送之后删除缓存项。
	if err := client.rdb.Del(ctx, client.RedisCacheKey()).Err(); err != nil {
		log.Error("Del error: ", err)
		return err
	}

	client.currentSize = 0
	client.currentNum = 0
	// client.m.Unlock()
	return nil
}
