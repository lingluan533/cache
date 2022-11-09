package main

import (
	"cache/backend"
	"cache/dataStruct"
	"cache/iot_server"
	"cache/util/config"
	logger "cache/util/log"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	_ "net/http/pprof"
	"os"
	_ "runtime/pprof"
	"strconv"
	"time"
)

// flag 包来处理命令行参数,返回值都是指针类型
var node = flag.String("node", "redis1", "节点名称参数")
var status = flag.Bool("status", false, "状态监测参数")

//var systemStatus backend.SystemStatus

func main() {
	fmt.Println("系统启动")
	flag.Parse() //解析flag，把用户传递的命令行参数解析为对应变量的值
	logger.Init()
	ctx := context.Background()
	//初始化etcd客户端
	config.GetETCDClient()
	config := config.Initialize()

	log.Println("配置文件加载完成")
	log.Println(config)
	//读取config.yaml中的内容
	if !config.Common.LogConfig.OutputFile {
		log.SetOutput(os.Stdout)
	}
	//交易信息打包推送到hraft节点
	//go backend.StartGrpcPort(":" + "8880")
	redisClient, err := backend.NewRedisBackend(ctx, &config)

	if err != nil {
		log.Error("后台程序运行失败: ", err)
		panic(err)
	}
	//status设置为TRUE的时候启动状态检测，RedisClient就是rdb，一个客户端的redis连接
	if *status {
		log.Info("启动系统状态监测...")
		log.Info("当前存活节点：", *node)
		stringCmd := redisClient.Info(ctx, "keyspace")
		log.Info("缓存空间状态：", stringCmd)
		log.Info("接收到终端的请求数量：", backend.Status.RequestNumber)
		log.Info("时间触发次数：", backend.Status.TimePushNumber)
		log.Info("数据量触发次数：", backend.Status.SizePushNumber)

	}
	//通道results是IoT节点传来的数据
	results := make(chan interface{})
	defer close(results)

	go consumeResult(ctx, redisClient, &config, results)

	ticker := time.NewTicker(time.Second * 10)

	go func() {
		n, err := redisClient.Exists(ctx, "ReceiptSet").Result()
		if err != nil {
			log.Error("quaryReceipt error: ", err)
		}
		m, err := redisClient.Exists(ctx, "TransactionSet").Result()
		if err != nil {
			log.Error("quaryTranscation error: ", err)
		}
		timestamp := time.Now().Unix()
		for {
			<-ticker.C

			if n > 0 {
				//清理存证过期项
				tm := time.Unix(timestamp, 0)
				//设置清理的时间范围
				m, _ := time.ParseDuration("-168h")
				//一周前
				tm2 := tm.Add(m)

				res2, err := redisClient.ZRemRangeByScore(ctx, "ReceiptSet", "0", strconv.FormatInt(tm2.Unix(), 10)).Result()
				if err != nil {
					log.Error("receiptRem error: ", err)
				}
				log.Info("Receipt Rem Success: ", res2)

			}

			if m > 0 {
				//清理交易过期项
				tm := time.Unix(timestamp, 0)
				m, _ := time.ParseDuration("-168h")
				//一周前
				tm2 := tm.Add(m)

				res2, err := redisClient.ZRemRangeByScore(ctx, "TransactionSet", "0", strconv.FormatInt(tm2.Unix(), 10)).Result()
				if err != nil {
					log.Error("transactionRem error: ", err)
				}
				log.Info("Transaction Rem Success: ", res2)
			}
			//隔半天清理一次
			ticker.Reset(time.Minute * 720)
		}
	}()
	go Monitor(ctx, redisClient)
	httpServer := iot_server.NewIOTServer(ctx, results, redisClient)
	// consul client registe
	consulClientRegiste(&config)
	//启动echo服务
	httpServer.Logger.Fatal(httpServer.Start(config.Cache.RedisGroup["redis1"].WebService.URL))
}

//iot节点交易信息缓存
func consumeResult(ctx context.Context, rdb *redis.Client, config *dataStruct.GlobalConfig, results chan interface{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case result := <-results:
			switch v := result.(type) {
			//question:这里不能用全局变量？ 存证数据记录
			case *iot_server.DataReceipt:
				//:TODO 所有分支都执行的关键字是什么来着
				backend.Status.RequestNumber++
				data, err := json.Marshal(v)
				if err != nil {
					log.Error("marshal error: ", err)
					continue
				}
				//发布订阅
				if err := rdb.Publish(ctx, backend.LedgerMap[v.DataType], string(data)).Err(); err != nil {
					log.Error("publish error: ", err)
					return
				}
				//				log.Info("Publish Success: ", string(data))
				//redis set存入值
				//start := time.Now().UnixNano()
				statusCmd := rdb.Set(ctx, v.KeyId, string(data), time.Second*time.Duration(config.Cache.CommonConfig.ExpireTime))
				if statusCmd == nil || statusCmd.Err() != nil {
					log.Error("Set error: ", statusCmd.Err())
					return
				}
				//				log.Info("Set Success: ", string(data))
				timestamp := time.Now().Unix()
				rdb.ZAdd(ctx, "ReceiptSet", &redis.Z{
					Score:  float64(timestamp),
					Member: string(data),
				})
				//zset的score存储时间值,每天定时扫描一下哪些过期的,取当前时间七天前的时间戳值，然后遍历zset找到比这个小的值都删除掉，就删除了七天前的数据

				//				log.Info("Receipt sSet Success: ", string(data))
			//实时交易记录
			case *iot_server.DataTransaction:
				backend.Status.RequestNumber++
				data, err := json.Marshal(v)
				if err != nil {
					log.Error("marshal error: ", err)
					continue
				}
				if err := rdb.Publish(ctx, backend.LedgerMap[v.DataType], string(data)).Err(); err != nil {
					log.Error("publish error: ", err)
					return
				}
				//log.Info("Publish Success: ", string(data))
				statusCmd := rdb.Set(ctx, v.TransactionId, string(data), time.Second*time.Duration(config.Cache.CommonConfig.ExpireTime))
				if statusCmd == nil || statusCmd.Err() != nil {
					log.Error("Set error: ", statusCmd.Err())
					return
				}
				//log.Info("Set Success: ", string(data))
				timestamp := time.Now().Unix()
				rdb.ZAdd(ctx, "TransactionSet", &redis.Z{
					Score:  float64(timestamp),
					Member: string(data),
				})

				log.Info("TransactionSet zSet Success: ", string(data))
			default:
				log.Error("error type %+v", v)
			}
		}
	}
}
func Monitor(ctx context.Context, redisClient *redis.Client) {
	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		timestamp := time.Now().Unix()
		log.Info("开始更新zSet:", timestamp)
		backend.GetValue(ctx, redisClient, timestamp)
		//隔半分钟查询一次
		ticker.Reset(time.Second * 30)
	}
}

func consulClientRegiste(globalConfig *dataStruct.GlobalConfig) {
	// 服务注册
	// 1.初始化 consul 配置
	consulConfig := api.DefaultConfig()
	config := globalConfig.Consul
	consulConfig.Address = config.ConsulAddress + ":" + config.ConsulPort
	// 2.创建 consul 对象
	consulClient, _ := api.NewClient(consulConfig)
	// 3.注册的服务配置
	log.Println(config.Name)
	reg := api.AgentServiceRegistration{
		ID:                config.ID,
		Name:              config.Name,
		Tags:              []string{"consul", "http", "edgenode"},
		Address:           config.LocalAddress,
		Port:              config.LocalServicePort,
		EnableTagOverride: true,
		Check: &api.AgentServiceCheck{
			CheckID:  config.HealthCheckID + config.LocalAddress,
			TCP:      config.HealthTCP,
			Timeout:  config.HealthTimeout,
			Interval: config.HealthInterval,
		},
	}
	log.Println(reg)
	log.Println("http://" + config.LocalAddress + ":" + strconv.Itoa(config.LocalServicePort) + "/health")
	// 4. 注册 grpc 服务到 consul 上
	err := consulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Println(err)
	}
}
