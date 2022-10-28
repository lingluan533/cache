package backend

import (
	pb "cache/backend/rpc/proto"
	"cache/dataStruct"
	"context"
	"encoding/json"
	_ "fmt"

	clientv3 "go.etcd.io/etcd/client/v3"

	"math/rand"
	_ "net/rpc"
	_ "sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//var ServiceDiscoveryFlag bool = false

func NewRpcClient(ctx context.Context, config *dataStruct.GlobalConfig, ledgerName string, etcdClient *clientv3.Client, results []string) error {
	// set up a connection to the server.
	connCtx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.Cache.CommonConfig.Connection))
	defer cancel()
	//与hraft的ip地址建立grpc连接
	conn, err := grpc.DialContext(connCtx, config.Consensus.EtcdGroup[config.Common.LedgerName[ledgerName].Leader].HraftGrpcAddress, grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		// :todo 容错机制
		log.Error("did not connect to the leader in config.yaml: ", config.Consensus.EtcdGroup[config.Common.LedgerName[ledgerName].Leader].HraftGrpcAddress)

		// 节点发现与轮询机制
		log.Info("start ServiceDiscovery...")
		nodeList := ServiceDiscovery(ctx, etcdClient, ledgerName)
		if nodeList == nil {
			log.Info("no active node was found")
			return err
		}
		log.Info("find active node list: ", nodeList)
		nodeFound := RoundRobin(nodeList, ledgerName)
		log.Info("try to connect: ", nodeFound)

		// 第二次尝试连接
		connCtx2, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.Cache.CommonConfig.Connection))
		defer cancel()

		conn, err = grpc.DialContext(connCtx2, nodeFound, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Error("did not connect: ", nodeFound)
			return err
		}
	}

	defer conn.Close()

	c := pb.NewToUpperClient(conn)

	subCtx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.Cache.CommonConfig.Response))
	defer cancel()

	// :todo 发记录数组
	switch ledgerName {
	case LedgerMap[DataTypeNodeCredibility]:
		var nodeCredibilityBlock pb.NodeCredibleData
		for _, result := range results {
			var nodeCredibilityTransaction pb.Transaction
			if err := json.Unmarshal([]byte(result), &nodeCredibilityTransaction); err != nil {
				log.Error("unmarshal error: ", err)
				return err
			}
			nodeCredibilityBlock.Transactions = append(nodeCredibilityBlock.Transactions, &nodeCredibilityTransaction)
		}
		//log.Info("打包好的数据: ", nodeCredibilityBlock.Transactions)
		r, err := c.NodeCredible(subCtx, &pb.NodeCredibleData{Transactions: nodeCredibilityBlock.Transactions})
		if err != nil {
			log.Error("could not greet: ", err)
			return err
		}
		log.Info("返回信息: ", r)
	case LedgerMap[DataTypeVideo]:
		var videoBlock pb.VideoData
		for _, result := range results {
			var videoReceipt pb.DataReceipt
			if err := json.Unmarshal([]byte(result), &videoReceipt); err != nil {
				log.Error("unmarshal error: ", err)
				return err
			}
			videoBlock.DataReceipts = append(videoBlock.DataReceipts, &videoReceipt)
		}
		//log.Info("打包好的数据: ", videoBlock.DataReceipts)
		r, err := c.Video(subCtx, &pb.VideoData{DataReceipts: videoBlock.DataReceipts})
		if err != nil {
			log.Error("cloud not greet: ", err)
			return err
		}
		log.Infof("返回信息: %s", r)
	case LedgerMap[DataTypeSensor]:
		var sensorBlock pb.SensorData
		for _, result := range results {
			var sensorTransaction pb.Transaction
			if err := json.Unmarshal([]byte(result), &sensorTransaction); err != nil {
				log.Error("unmarshal error: ", err)
				return err
			}
			sensorBlock.Transactions = append(sensorBlock.Transactions, &sensorTransaction)
		}
		//log.Info("打包好的数据:", sensorBlock.Transactions)
		r, err := c.Sensor(subCtx, &pb.SensorData{Transactions: sensorBlock.Transactions})
		if err != nil {
			log.Error("cloud not greet: ", err)
			return err
		}
		log.Info("返回信息: ", r)
	case LedgerMap[DataTypeUserBehaviour]:
		var userBehaviourBlock pb.UserBehaviourData
		for _, result := range results {
			var userBehaviourReceipt pb.DataReceipt
			if err := json.Unmarshal([]byte(result), &userBehaviourReceipt); err != nil {
				log.Error("unmarshal error: ", err)
				return err
			}
			userBehaviourBlock.DataReceipts = append(userBehaviourBlock.DataReceipts, &userBehaviourReceipt)
		}
		//log.Info("打包好的数据: ", userBehaviourBlock.DataReceipts)
		r, err := c.UserBehaviour(subCtx, &pb.UserBehaviourData{DataReceipts: userBehaviourBlock.DataReceipts})
		if err != nil {
			log.Error("cloud not greet: ", err)
			return err
		}
		log.Info("返回信息: ", r)
	case LedgerMap[DataTypeAccessLog]:
		var serviceAccessBlock pb.ServiceAccessData
		for _, result := range results {
			var serviceAccessTransaction pb.Transaction
			if err := json.Unmarshal([]byte(result), &serviceAccessTransaction); err != nil {
				log.Error("unmarshal error: ", err)
				return err
			}
			//log.Info("Send Transaction: ", result)
			serviceAccessBlock.Transactions = append(serviceAccessBlock.Transactions, &serviceAccessTransaction)
		}
		//log.Info("打包好的数据: ", serviceAccessBlock.Transactions)
		r, err := c.ServiceAccess(subCtx, &pb.ServiceAccessData{Transactions: serviceAccessBlock.Transactions})
		if err != nil {
			log.Error("cloud not greet: ", err)
			return err
		}
		log.Info("返回信息: ", r)
	}

	return nil
}

func RoundRobin(nodeList []string, ledgerName string) string {
	nodeFound := nodeList[rand.Intn(len(nodeList))]
	return nodeFound
}
