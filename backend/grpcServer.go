package backend

import (
	pb "cache/backend/rpc/proto"
	"net"
	"strings"
	"sync"

	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
)

type server struct{} //定义一个server结构体

// func RedisServerMain(client *clientv3.Client) { //redis服务端
// 	//clientRedis := client
// 	//遍历开启端口
// 	go StartGrpcPort(":" + "8880")
// }

//遍历开启端口
func StartGrpcPort(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Error("开启端口失败: ", err)
	}
	log.Info("端口开启成功！", Port)

	s := grpc.NewServer()
	pb.RegisterToUpperServer(s, &server{})
	//reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Error("端口服务调用失败: ", err)
	}
}

//grpc测试接口
func (s *server) Upper(ctx context.Context, in *pb.UpperRequest) (*pb.UpperReply, error) {
	//log.Info("测试rpc，接收到数据:", in.Name)
	return &pb.UpperReply{Message: strings.ToUpper(in.Name)}, nil
}

//Video账本类型  存证数据
func (s *server) Video(ctx context.Context, in *pb.VideoData) (*pb.Response, error) {
	//为存储数据的全局变量加锁
	TransactionDatamu = new(sync.RWMutex)
	ReceiptDatamu = new(sync.RWMutex)
	MDDatamu = new(sync.RWMutex)
	//log.Info("Video账本类型，接收到数据: ", in.DataReceipts)
	return &pb.Response{ErrCode: SuccessCode, ErrMsg: ""}, nil
}

//UserBehaviour账本类型  存证数据
func (s *server) UserBehaviour(ctx context.Context, in *pb.UserBehaviourData) (*pb.Response, error) {
	//为存储数据的全局变量加锁
	TransactionDatamu = new(sync.RWMutex)
	ReceiptDatamu = new(sync.RWMutex)
	MDDatamu = new(sync.RWMutex)
	//log.Info("UserBehaviour账本类型，接收到数据:", in.DataReceipts)
	return &pb.Response{ErrCode: SuccessCode, ErrMsg: ""}, nil
}

//NodeCredible账本类型  交易数据
func (s *server) NodeCredible(ctx context.Context, in *pb.NodeCredibleData) (*pb.Response, error) {
	//为存储数据的全局变量加锁
	TransactionDatamu = new(sync.RWMutex)
	ReceiptDatamu = new(sync.RWMutex)
	MDDatamu = new(sync.RWMutex)
	//log.Info("NodeCredible账本类型，接收到数据:", in.Transactions)
	return &pb.Response{ErrCode: SuccessCode, ErrMsg: ""}, nil
}

//Sensor账本类型  交易数据
func (s *server) Sensor(ctx context.Context, in *pb.SensorData) (*pb.Response, error) {
	//为存储数据的全局变量加锁
	TransactionDatamu = new(sync.RWMutex)
	ReceiptDatamu = new(sync.RWMutex)
	MDDatamu = new(sync.RWMutex)
	//log.Info("Sensor账本类型，接收到数据:", in.Transactions)
	return &pb.Response{ErrCode: SuccessCode, ErrMsg: ""}, nil
}

//ServiceAccess账本类型  交易数据
func (s *server) ServiceAccess(ctx context.Context, in *pb.ServiceAccessData) (*pb.Response, error) {
	//为存储数据的全局变量加锁
	TransactionDatamu = new(sync.RWMutex)
	ReceiptDatamu = new(sync.RWMutex)
	MDDatamu = new(sync.RWMutex)
	//log.Info("ServiceAccess账本类型，接收到数据:", in.Transactions)
	return &pb.Response{ErrCode: SuccessCode, ErrMsg: ""}, nil
}
