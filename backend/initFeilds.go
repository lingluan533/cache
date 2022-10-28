package backend

import (

	//"hraft/utils"

	"sync"
	"time"
)

var (
	BLOCK_TYPE_MIN      = "MINUTE"
	BLOCK_TYPE_TENMINUT = "TENMINUTE"
	BLOCK_TYPE_DAY      = "DAY"
	ErrMsg              = "未知错误"
	SuccessCode         = int32(200)
	ErrCode             = int32(400)
	//用于拼接各字段值 之间分割符
	KeySplit = ":"
	//用于分割 时间戳 和 KeyId 或者 TransactionId
	TIMESTAMP_KEYID = "#"

	LEDGER_TYPE_VIDEO         = "video"
	LEDGER_TYPE_USER_BEHAVIOR = "user_behaviour"

	LEDGER_TYPE_NODE_CREDIBLE  = "node_credible"
	LEDGER_TYPE_SENSOR         = "sensor"
	LEDGER_TYPE_SERVICE_ACCESS = "service_access"

	ALL_LEDGER_TYPE_ARRAY = []string{LEDGER_TYPE_VIDEO, LEDGER_TYPE_USER_BEHAVIOR, LEDGER_TYPE_NODE_CREDIBLE, LEDGER_TYPE_SENSOR, LEDGER_TYPE_SERVICE_ACCESS}
	BLOCK_TYPE_ARRAY      = []string{BLOCK_TYPE_MIN, BLOCK_TYPE_TENMINUT, BLOCK_TYPE_DAY}

	//现在该结点启动账本数组
	GlobalLedgerArray = []string{}

	//定义存储数据的map结构，数据不将存储于数据库
	TransactionData = make(map[string]string)
	ReceiptData     = make(map[string]string)
	MDData          = make(map[string][]string)
	//ReceiptMDData     = make(map[string][]string)
	//为map加锁变量
	TransactionDatamu *sync.RWMutex
	ReceiptDatamu     *sync.RWMutex
	MDDatamu          *sync.RWMutex

	TenMinBlockChangeVideo         = make([]bool, 144)
	TenMinBlockChangeUserBehavior  = make([]bool, 144)
	TenMinBlockChangeNodeCredible  = make([]bool, 144)
	TenMinBlockChangeSensor        = make([]bool, 144)
	TenMinBlockChangeServiceAccess = make([]bool, 144)

	DailyChangeVideo         = false
	DailyChangeUserBehavior  = false
	DailyChangeNodeCredible  = false
	DailyChangeSensor        = false
	DailyChangeServiceAccess = false
)

var DialTimeout time.Duration
var RequestTimeout time.Duration
var Port string

var GlobalLeaderId uint64
var GlobalLeaderName string

// DialTimeout = time.Duration(10) * time.Second
// RequestTimeout = time.Duration(10) * time.Second

// Port = strings.Split(config.Consensus.EtcdGroup[GlobalLeaderName].HraftAddress, ":")[1]
