package backend

import (
	"os"
	"context"

	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	_"time"
	"strconv"
	"strings"
	log "github.com/sirupsen/logrus"
)

//type BlockHeaders struct{
//	Header []BlockHeader `json: "Header"`
//}
type MinuteDataBlock struct {
	Header       BlockHeader   `json:"Header,omitempty"`
	DataReceipts []DataReceipt `json:"DataReceipts,omitempty"` //元数据
}
type BlockHeader struct {
	CreateTimestamp string `json:"CreateTimestamp,omitempty"` //创建时间戳
	KeyId           string `json:"keyId,omitempty"`
	BlockHeight     int64  `json:"BlockHeight,omitempty"` //通过该字段，获取当前区块 可以使用不同链
	//具体数据结构类型
	DataType         string `json:"DataType,omitempty"`                   //数据类型
	DataValue        string `json:"DataValue,omitempty"`                 //数据价值
	UpdateTimestamp  string `json:"UpdateTimestamp,omitempty"`     //更新时间戳
	DataHash         string `json:"DataHash,omitempty"`                   //数据哈希值
	BlockHash        string `json:"BlockHash,omitempty"`                 //区块哈希值
	PreBlockHash     string `json:"PreBlockHash,omitempty"`           //前一个区块hash值
	Nonce            int32  `json:"Nonce,omitempty"`                       //nonce 值
	Target           int32  `json:"Target,omitempty"`                     //目标值
	CurrentDataCount int64  `json:"CurrentDataCount,omitempty"` //当前数据记录量
	CurrentDataSize  int64  `json:"CurrentDataSize,omitempty"`   //当前数据大小
	Version          string `json:"Version,omitempty"`                    //版本号
	BlockType        string `json:"BlockType,omitempty"`                //区块类型
	LedgerType       string `json:"LedgerType,omitempty"`              //账本类型
}
type DataReceipt struct {
	CreateTimestamp     string   `json:"createTimestamp"`
	EntityId            string   `json:"entityId"`
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
type DataReceiptBlockInfo struct {
	DataReceipt
	BlockHeight     int64  `json:"BlockHeight"` //通过该字段，获取当前区块 可以使用不同链
	BlockHash        string `json:"BlockHash"`                 //区块哈希值
	CurrentDataCount int64  `json:"CurrentDataCount"` //当前数据记录量
	CurrentDataSize  int64  `json:"CurrentDataSize"`   //当前数据大小
	Version          string `json:"Version"`                    //版本号
	BlockType        string `json:"BlockType"`                //区块类型
}
// 实时交易记录
type Transaction struct {
	Timestamp     string  `json:"timeStamp" validate:"required"`
	EntityId      string  `json:"entityId"`
	TransactionId string  `json:"transactionId" validate:"required"`
	Initiator     string  `json:"initiator"`
	Recipient     string  `json:"recipient"`
	TxAmount      float64 `json:"txAmount"`
	DataType      string  `json:"dataType" validate:"required"`
	ServiceType   string  `json:"serviceType"`
	Remark        string  `json:"remark"`
}
//返回读取到的文件信息
func ReadTxMinFiletoTenmin(time string, ledger string, index string) MinuteDataBlock{
	//time是日期，index是当前分钟数
	//mindirectory, err := os.Getwd()
	//if err != nil {
	//	log.Error("获取当前路径失败 =: ", err)
	//}
	var (
		fileName = "D:\\Go\\src\\hraft1102" + "/scope/" + time + "/" + ledger + "/MINUTE" + "/" + index
	//var (
	//	fileName = "E:\\Go\\go\\src\\cache" + "/scope/" + time + "/" + ledger + "/MINUTE" + "/" + index
	//)
	//fmt.Println(fileName)
	log.Info("区块文件名：", fileName)
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Info("区块文件不存在")
		return MinuteDataBlock{}
	}else {
		log.Info("文件存在，开始更新")
		jsonfile,_ := os.Open(fileName)
		defer jsonfile.Close()
		//读取文件
		fileContent, err := ioutil.ReadAll(jsonfile)
		if err != nil {
			log.Error("Read file err =: ", err)
		}
		//fmt.Println(string(fileContent))
		var minBlock MinuteDataBlock
		//var user User
		if err := json.Unmarshal([]byte(fileContent), &minBlock); err != nil {
			log.Error("反解析 file error =: ", err)
		}
		//fmt.Println(minBlock.Header)
		//fmt.Println(len(minBlock.DataReceipts))
		return minBlock
	}
}

func Convert(blockInfo BlockHeader, receipt DataReceipt) *DataReceiptBlockInfo{
	return &DataReceiptBlockInfo{
		receipt,
		blockInfo.BlockHeight,
		blockInfo.BlockHash,
		blockInfo.CurrentDataCount,
		blockInfo.CurrentDataSize,
		blockInfo.Version,
		blockInfo.BlockType,
	}
}
//2022-07-12 10:06:00.0450021
//根据redis存证数据的时间戳返回分钟数
func GetIndexMinInt(x string) (string, int) {
	//timeFormatString := time.Now().Format("2006-01-02 15:04:05")
	//var naosecond = time.Now().Nanosecond()/1e6
	dayTimeArray := strings.Split(x, " ")
	minTimeArray := strings.Split(dayTimeArray[1], ":")
	hourInt, _ := strconv.Atoi(minTimeArray[0])
	minInt, _ := strconv.Atoi(minTimeArray[1])
	//最后返回的是当日0点开始算至今过去的分钟数
	indexMinInt := hourInt*60 + minInt
	//indexMinString := strconv.Itoa(indexMinInt)
	//timeCorrect:=fmt.Sprintf("%s.%d",timeFormatString,naosecond)
	return dayTimeArray[0], indexMinInt
}
//查询redis中的所有Key，根据Key的存储时间和类型到hraft文件中查找到该key，然后修改这个key在zset中的值，添加区块属性
//每分钟定时查询一次上一分钟范围内的新加元素，如果产生了新元素，则去该分钟对应的文件中查找区块信息并更新。
func GetValue(ctx context.Context, rdb *redis.Client, timeStamp int64){
	//timestamp := time.Now().Unix()
	//timestamp := 1657591590
	left := float64(timeStamp)-100.0
	right := float64(timeStamp)
	n, err := rdb.ZCount(ctx, "ReceiptSet", fmt.Sprint(left), fmt.Sprint(right)).Result()
	if err != nil {
		log.Error("zcount查询失败 : ", err)
	}

	if n > 0{
		log.Info("当前zSet中被选中数据个数为：",n )
		res, err := rdb.ZRangeByScore(ctx, "ReceiptSet", &redis.ZRangeBy{
			Min:    fmt.Sprint(left),
			Max:    fmt.Sprint(right),
		}).Result()
		if err != nil {
			log.Error("zrange查询失败 : ", err)
		}
		//fmt.Println(res)
		//val, err := rdb.Get(ctx, keyId).Result()
		var result []DataReceipt
		for _,n := range res{

			var tmp DataReceipt
			if err := json.Unmarshal([]byte(n), &tmp); err != nil {
				log.Error("反解析 file : ", err)
			}
			//log.Info("开始更新zSet")
			result = append(result, tmp)
		}
		for _,n := range result{
			//log.Info("开始更新zSet")
			//fmt.Println(n.ServiceType, " ", n.CreateTimestamp)
			time1, index := GetIndexMinInt(n.CreateTimestamp)
			//根据类型和时间，查询对应区块文件，将区块信息加入set中
			log.Info("分钟块号码为：", index)
			blockInfo := ReadTxMinFiletoTenmin(time1, n.ServiceType, strconv.Itoa(index))
			if blockInfo.Header.KeyId != "" {
				res1, _ := json.Marshal(n)
				//fmt.Println(string(res1))
				score,_ := rdb.ZScore(ctx, "ReceiptSet", string(res1)).Result()
				if score == 0.0{
					continue
				}
				log.Info("当前更新的值的score：", score)
				//fmt.Println(score)
				res, err := rdb.ZRem(ctx, "ReceiptSet", string(res1)).Result()
				if err != nil {
					log.Error("删除失败 error: ", err)
				}
				log.Info("删除zSet中数据个数：", res)
				//score.Result()
				end := Convert(blockInfo.Header, n)
				end2, _ := json.Marshal(*end)
				rdb.ZAdd(ctx, "ReceiptSet", &redis.Z{
					Score:  score,
					Member: string(end2),
				})
			}else {
				break
			}
		}
	}else {
		log.Info("当前zSet中没有新数据" )
	}

}

