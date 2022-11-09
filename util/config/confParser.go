package config

import (
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
	"net"
	"strconv"
	"strings"
	"time"

	//yaml 包使 Go 程序能够轻松地对 YAML 值进行编码和解码。
	"cache/backend"
	"cache/dataStruct"
	"io/ioutil"
)

var EtcdClient *clientv3.Client
var GlobalConfig dataStruct.GlobalConfig
var DialTimeout time.Duration
var RequestTimeout time.Duration

func Initialize() dataStruct.GlobalConfig {
	data, err := ioutil.ReadFile("./config.yaml")
	//ReadFile 从filename指定的文件中读取数据并返回文件的内容。
	if err != nil {
		log.Fatal(err)
	}

	var config dataStruct.GlobalConfig
	//加载配置文件data，config获取全部配置文件中的内容。
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Error("配置文件加载错误: ", err)
		log.Fatal(err)
	}
	//???
	backend.LedgerMap[backend.DataTypeNodeCredibility] = "node_credible"
	backend.LedgerMap[backend.DataTypeVideo] = "video"
	backend.LedgerMap[backend.DataTypeSensor] = "sensor"
	backend.LedgerMap[backend.DataTypeUserBehaviour] = "user_behaviour"
	backend.LedgerMap[backend.DataTypeAccessLog] = "service_access"
	ip := Ips()
	config.Consul.LocalAddress = ip
	config.Consul.HealthTCP = ip + ":" + strconv.Itoa(config.Consul.LocalServicePort)
	config.Consul.ID = config.Consul.Name + "_" + ip
	GlobalConfig = config
	//yaml文件中配置了10秒
	DialTimeout = time.Duration(GlobalConfig.Consensus.CommonConfig.Timeout) * time.Second
	//这里设置了RequestTimeout的值，后面很多地方都用到了这个值
	RequestTimeout = time.Duration(GlobalConfig.Consensus.CommonConfig.Timeout) * time.Second

	return config
}

func Ips() string {

	ips := make(map[string]string)

	interfaces, _ := net.Interfaces()

	for _, i := range interfaces {
		byName, _ := net.InterfaceByName(i.Name)
		addresses, _ := byName.Addrs()
		for _, v := range addresses {
			ips[byName.Name] = v.String()
			fmt.Println(byName.Name, v.String(), v.Network())
			if strings.HasPrefix(v.String(), "192.168.195.") {
				fmt.Println("检测到ip:", v.String())
				return strings.TrimSuffix(v.String(), "/24")
			}
		}
		for _, v := range addresses {
			ips[byName.Name] = v.String()
			//fmt.Println(byName.Name, v.String(), v.Network())
			if strings.HasPrefix(v.String(), "192.168.216.") {
				fmt.Println("检测到ip:", v.String())
				return strings.TrimSuffix(v.String(), "/24")
			}
		}
	}
	return "127.0.0.1"

}

func GetETCDClient() {
	// 连接ETCD
	// etcd客户端
	var err error
	log.Printf("ETCD客户端连接中....")
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Error("etcd获取客户端错误", err)
		//return c.JSON(http.StatusInternalServerError, NewResult(err.Error(), nil))
	} else {
		status, _ := EtcdClient.Status(context.TODO(), "127.0.0.1:2379")
		log.Infof("etcd获取客户端成功：", EtcdClient)
		log.Infof("etcd客户端状态：", status)
	}
	log.Printf("ETCD客户端连接成功....")

}
