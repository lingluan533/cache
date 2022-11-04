package consul_service

import (
	"cache/dataStruct"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

func QueryOnlineEdgeNodes() []*consulapi.ServiceEntry {
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = service.ConsulAddress + ":" + service.ConsulPort
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		fmt.Println("consul client error : ", err)
		return nil
	}

	// 获取指定service
	serviceHealthy, _, err := client.Health().Service(service.Name, "", true, nil)
	if err != nil {
		fmt.Println("consul client error : ", err)
		return nil
	}
	return serviceHealthy
}

//TODO:暂时写死
var service = &dataStruct.ConsulConfig{
	ConsulAddress: "101.43.138.160",
	ConsulPort:    "8500",
	Name:          "EdgeNode",
}

func GetOneOnlineAddress() *consulapi.AgentService {
	//config := GetConfig()
	// 创建连接consul服务配置
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = service.ConsulAddress + ":" + service.ConsulPort
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		fmt.Println("consul client error : ", err)
		return nil
	}

	// 获取指定service
	serviceHealthy, _, err := client.Health().Service(service.Name, "", true, nil)
	if err != nil {
		fmt.Println("consul client error : ", err)
		return nil
	}
	return serviceHealthy[0].Service

}
