package backend

import "sync"

const (
	// DataReceipt
	DataTypeVideo         =  "video"
	DataTypeUserBehaviour = "user_behaviour"

	// DataTransaction
	DataTypeNodeCredibility = "node_credible"
	DataTypeSensor          = "sensor"
	DataTypeAccessLog       = "service_access"
)

var LedgerMap = make(map[string]string)

type SystemStatus struct {
	RecordMutex sync.Mutex
	RequestNumber int
	SizePushNumber int
	TimePushNumber int
	BlockInQueue int
	FailPushNumber int
}
var Status SystemStatus


