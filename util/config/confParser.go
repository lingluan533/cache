package config

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"

	//yaml 包使 Go 程序能够轻松地对 YAML 值进行编码和解码。
	"cache/backend"
	"cache/dataStruct"
	"io/ioutil"
)

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

	return config
}
