package iot_server

// Result represents HTTP response body.
type Result struct {
	Err  interface{} `json:"err"`  // error message
	Data interface{} `json:"data"` // data object
}

// NewResult creates a result with Code=0, Msg="", Data=nil.
func NewResult(params ...interface{}) *Result {
	paramLen := len(params)
	if paramLen == 0 {
		return &Result{
			Err:  nil,
			Data: nil,
		}
	} else if paramLen == 1 {
		return &Result{
			Err:  nil,
			Data: params[0],
		}
	} else {
		return &Result{
			Err:  params[0],
			Data: params[1],
		}
	}
}

type Receipt struct {
	KeyId               string   `json:"keyId" validate:"required"`
	ReceiptValue        float64  `json:"receiptValue"`
	Version             string   `json:"version"`
	UserName            string   `json:"userName"`
	OperationType       string   `json:"operationType"`
	DataType            string   `json:"dataType" validate:"required"`
	ServiceType         string   `json:"serviceType"`
	FileName            string   `json:"fileName"`
	FileSize            float64  `json:"fileSize"`
	FileHash            string   `json:"fileHash"`
	Uri                 string   `json:"uri"`
	ParentKeyId         string   `json:"parentKeyId"`
	AttachmentFileUris  []string `json:"attachmentFileUris"`
	AttachmentTotalHash string   `json:"attachmentTotalHash"`
}

// 存证数据记录
type DataReceipt struct {
	CreateTimestamp string `json:"createTimestamp" validate:"required"`
	EntityId  string `json:"entityId"`

	Receipt //`json:"receipt"`
}

type DataReceipts struct {
	CreateTimestamp  string    `json:"createTimestamp" validate:"required"`
	EntityId   string    `json:"entityId"`
	DataValue  float64   `json:"dataValue"`
	DataRecNum int64     `json:"dataRecNum"`
	Receipts   []Receipt `json:receipts`
}

type BlockInfoResp struct {
	KeyId       string `json:"keyId"`
	TxId        string `json:"txId"`
	BlockHeight int64  `json:"blockHeight"`
	BlockHash   string `json:"blockHash"`
}

type ReceiptResponseInfo struct {
	Success       bool          `json:"success"`
	Status        bool          `json:"status"`
	Err           string        `json:"err"`
	DataReceipt   DataReceipt   `json:"dataReceipt"`
	BlockInfoResp BlockInfoResp `json:"blockInfoResp"`
}

type TransactionResponseInfo struct {
	Success       bool          `json:"success"`
	Status        bool          `json:"status"`
	Err           string        `json:"err"`
	DataTransaction   DataTransaction   `json:"dataTransaction"`
	BlockInfoResp BlockInfoResp `json:"blockInfoResp"`
}

type Transaction struct {
	TransactionId string  `json:"transactionId" validate:"required"`
	Initiator     string  `json:"initiator"`
	Recipient     string  `json:"recipient"`
	TxAmount      float64 `json:"txAmount"`
	DataType      string  `json:"dataType" validate:"required"`
	ServiceType   string  `json:"serviceType"`
	Remark        string  `json:"remark"`
}

// 实时交易记录
type DataTransaction struct {
	CreateTimestamp string `json:"createTimestamp" validate:"required"`
	EntityId  string `json:"entityId"`

	Transaction
}

type DataTransactions struct {
	CreateTimestamp    string        `json:"createTimestamp" validate:"required"`
	EntityId     string        `json:"entityId"`
	TurnOver     float64       `json:"turnOver"`
	TxRecNum     int64         `json:"txRecNum"`
	Transactions []Transaction `json:"transactions"`
}

