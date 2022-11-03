package iot_server

import (
	"cache/backend"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var LOC, _ = time.LoadLocation("Asia/Shanghai")

func NewIOTServer(ctx context.Context, results chan interface{}, rdb *redis.Client) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = ErrorHandler
	//验证器
	e.Validator = &CustomValidator{validator: validator.New()}
	//e.Validator = validator.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Scope!")
	})
	i := 0
	e.GET("/health", func(c echo.Context) error {
		fmt.Println("Consul health check!!")
		return c.JSON(http.StatusOK, NewResult(nil, nil))
	})
	//TODO:登录请求,在区块链上查找用户信息  在etcd集群中找
	e.POST("/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "Login success!")
	})
	e.POST("/storeReceipt", func(c echo.Context) error {
		//log.Info("接收到存证数据")

		var receipts DataReceipts
		if err := c.Bind(&receipts); err != nil {
			// {"err": "marshal error", "data": nil}
			log.Error("marshal error: ", err)
			return c.JSON(http.StatusOK, NewResult(err.Error(), nil))
		}
		if err := c.Validate(&receipts); err != nil {
			log.Error("validate error: ", err)
			return c.JSON(http.StatusOK, NewResult(err.Error(), nil))
		}

		for _, r := range receipts.Receipts {
			var receipt DataReceipt
			var naosecond = time.Now().UnixNano()
			stringnase := strconv.FormatInt(naosecond, 10)
			nasecount := len(stringnase)
			naosecond1 := stringnase[nasecount-5 : nasecount-3]
			naosecond2 := stringnase[nasecount-3 : nasecount-1]
			naosecond3 := stringnase[nasecount-1 : nasecount]
			//fmt.Println(stringnase)
			timeUnix := time.Now().Unix() //时间戳
			//该时间传入的整个数组集都统一时间戳
			str := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
			receipt.Receipt = r
			timeCorrect := fmt.Sprintf("%s.%s%s%s", str, naosecond1, naosecond2, naosecond3)
			//fmt.Println(receipt.DataType, backend.DataTypeVideo)
			if receipt.DataType == backend.DataTypeVideo {
				timeCorrect = timeCorrect + "1"
			} else if receipt.DataType == backend.DataTypeUserBehaviour {
				timeCorrect = timeCorrect + "2"
			}
			//fmt.Println(timeCorrect)
			receipt.CreateTimestamp = timeCorrect
			receipt.EntityId = receipts.EntityId

			if err := c.Validate(&r); err != nil {
				log.Error("validate error: ", err)
				return c.JSON(http.StatusOK, NewResult(err.Error(), nil))
			}

			//result接收数据，将数据进行分类
			if receipt.DataType == backend.DataTypeVideo || receipt.DataType == backend.DataTypeUserBehaviour {
				i++
				if i%10000 == 0 {
					log.Infof("%v", i)
				}
				results <- &receipt
			} else {
				log.Error("DataType error")
				return c.JSON(http.StatusOK, NewResult("DataType error", nil))
			}
		}
		return c.JSON(http.StatusOK, NewResult("success"))
		//return c.JSON(http.StatusOK, receipts)
	})

	e.POST("/storeTransaction", func(c echo.Context) error {
		//log.Info("接收到交易数据")
		var transactions DataTransactions
		if err := c.Bind(&transactions); err != nil {
			log.Error("marshal error: ", err)
			return c.JSON(http.StatusOK, NewResult(err.Error(), nil))
		}
		if err := c.Validate(&transactions); err != nil {
			log.Error("validate error: ", err)
			return c.JSON(http.StatusOK, NewResult(err.Error(), nil))
		}
		// timeUnix := time.Now().Unix() //时间戳
		// //该时间传入的整个数组集都统一时间戳
		// str := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05.123")
		// transactions.CreateTimestamp = str[:23]
		//log.Info(transactions)
		for _, t := range transactions.Transactions {
			var transaction DataTransaction
			var naosecond = time.Now().UnixNano()
			stringnase := strconv.FormatInt(naosecond, 10)
			nasecount := len(stringnase)
			naosecond1 := stringnase[nasecount-5 : nasecount-3]
			naosecond2 := stringnase[nasecount-3 : nasecount-1]
			naosecond3 := stringnase[nasecount-1 : nasecount]
			//fmt.Println(stringnase)
			timeUnix := time.Now().Unix() //时间戳
			//该时间传入的整个数组集都统一时间戳
			str := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
			transaction.Transaction = t
			timeCorrect := fmt.Sprintf("%s.%s%s%s", str, naosecond1, naosecond2, naosecond3)
			//fmt.Println(timeCorrect)
			if transaction.DataType == backend.DataTypeSensor {
				timeCorrect = timeCorrect + "3"
			} else if transaction.DataType == backend.DataTypeNodeCredibility {
				timeCorrect = timeCorrect + "4"
			} else if transaction.DataType == backend.DataTypeAccessLog {
				timeCorrect = timeCorrect + "5"
			}

			transaction.CreateTimestamp = timeCorrect
			transaction.EntityId = transactions.EntityId

			if err := c.Validate(&t); err != nil {
				log.Error("validate error: ", err)
				return c.JSON(http.StatusOK, NewResult(err.Error(), nil))
			}
			if transaction.DataType == backend.DataTypeSensor || transaction.DataType == backend.DataTypeNodeCredibility || transaction.DataType == backend.DataTypeAccessLog {
				i++
				if i%10000 == 0 {
					log.Infof("%v", i)
				}
				results <- &transaction
			} else {
				log.Error("DataType error")
				return c.JSON(http.StatusOK, NewResult("DataType error", nil))
			}
		}
		return c.JSON(http.StatusOK, transactions)
	})
	//存证查询接口
	e.POST("/loadReceipt", func(c echo.Context) error {
		log.Info("查询存证数据")
		var respInfo ReceiptResponseInfo
		respInfo.DataReceipt.Receipt.KeyId = c.FormValue("KeyId")
		resp, err := rdb.Get(ctx, respInfo.DataReceipt.Receipt.KeyId).Result()

		if err != nil {
			respInfo.Success = false
			respInfo.Err = "GET error"
			log.Error("GET error: ", err)
			return c.JSON(http.StatusOK, resp)
		}
		if err := json.Unmarshal([]byte(resp), &respInfo.DataReceipt); err != nil {
			fmt.Println(string(resp))
			respInfo.Success = false
			respInfo.Err = "Unmarshal error"
			log.Error("Unmarshal error: ", err)
			return c.JSON(http.StatusOK, resp)
		}
		respInfo.Success = true
		respInfo.Status = true
		log.Info("查询存证数据成功")
		return c.JSON(http.StatusOK, respInfo)
	})
	//交易查询接口
	e.POST("/loadTransaction", func(c echo.Context) error {
		log.Info("查询交易数据")
		var respInfo TransactionResponseInfo
		respInfo.DataTransaction.Transaction.TransactionId = c.FormValue("KeyId")
		resp, err := rdb.Get(ctx, respInfo.DataTransaction.Transaction.TransactionId).Result()

		if err != nil {
			respInfo.Success = false
			respInfo.Err = "GET error"
			log.Error("GET error: ", err)
			return c.JSON(http.StatusOK, resp)
		}
		if err := json.Unmarshal([]byte(resp), &respInfo.DataTransaction); err != nil {
			fmt.Println(string(resp))
			respInfo.Success = false
			respInfo.Err = "Unmarshal error"
			log.Error("Unmarshal error: ", err)
			return c.JSON(http.StatusOK, resp)
		}
		respInfo.Success = true
		respInfo.Status = true
		log.Info("查询交易数据成功")
		return c.JSON(http.StatusOK, respInfo)
	})
	//e.POST("/queryTimeReceipt", func(c echo.Context) error {
	//	//c是提交的参数，
	//	// if err := c.Bind(&m); err != nil {
	//	// 	return c.JSON(http.StatusOK, NewResult(err.Error(), nil))
	//	// }
	//	log.Info("按时间戳查询存证数据")
	//	//var respInfo ReceiptResponseInfo
	//	//start := c.FormValue("StartTime")
	//	//end := c.FormValue("EndTime")
	//	start, _ := time.Parse("2006-01-02 15:04:05", c.FormValue("StartTime"))
	//	end, _ := time.Parse("2006-01-02 15:04:05", c.FormValue("EndTime"))
	//	res, err := rdb.ZRangeByScore(ctx, "zSet", &redis.ZRangeBy{
	//		Min: strconv.FormatInt(start.Unix(), 10),
	//		Max: strconv.FormatInt(end.Unix(), 10),
	//	}).Result()
	//
	//	if err != nil {
	//		// respInfo.Success = false
	//		// respInfo.Err = "GET error"
	//		log.Error("GET error: ", err)
	//		return c.JSON(http.StatusOK, res)
	//	}
	//	// fmt.Println(res)
	//	//将查到的多条结果赋值给数据结构
	//	// if err := json.Unmarshal([]byte(res), &respInfo.DataReceipt); err != nil {
	//	// 	fmt.Println(string(res))
	//	// 	respInfo.Success = false
	//	// 	respInfo.Err = "Unmarshal error"
	//	// 	log.Error("Unmarshal error: ", err)
	//	// 	return c.JSON(http.StatusOK, res)
	//	// }
	//	// respInfo.Success = true
	//	// respInfo.Status = true
	//	//log.Info("zSet Success: ", string(data))
	//	log.Info("查询存证数据成功: ", res)
	//	return c.JSON(http.StatusOK, res)
	//})
	type QueryBlocks struct {
		StartTime string `json:"StartTime" xml:"StartTime" form:"StartTime" query:"StartTime"`
		EndTime   string `json:"EndTime" xml:"EndTime" form:"EndTime" query:"EndTime"`
	}
	e.POST("/queryTimeReceipt", func(c echo.Context) error {

		st := c.FormValue("StartTime")
		et := c.FormValue("EndTime")
		log.Infof("按时间戳查询数据,startTime%v,endtime:%v", st, et)
		startTime, err := strconv.ParseInt(st, 10, 64)
		if err != nil {
			log.Error("queryTimeReceipt timestamp err")
			return c.JSON(http.StatusInternalServerError, errors.New("queryTimeReceipt startTime timestamp err"))
		}
		endTime, err := strconv.ParseInt(et, 10, 64)
		if err != nil {
			log.Error("queryTimeReceipt timestamp err")
			return c.JSON(http.StatusInternalServerError, errors.New("queryTimeReceipt endTime timestamp err"))
		}
		res, err := rdb.ZRangeByScore(ctx, "ReceiptSet", &redis.ZRangeBy{
			Min: strconv.FormatInt(startTime, 10),
			Max: strconv.FormatInt(endTime, 10),
		}).Result()
		fmt.Println(res)
		if err != nil {
			// respInfo.Success = false
			// respInfo.Err = "GET error"
			log.Error("GET error: ", err)
			return c.JSON(http.StatusOK, res)
		}
		log.Info("查询存证数据成功: ")
		return c.JSON(http.StatusOK, res)
	})
	e.POST("/queryTimeTransaction", func(c echo.Context) error {

		st := c.FormValue("StartTime")
		et := c.FormValue("EndTime")
		log.Infof("按时间戳查询交易数据,startTime%v,endtime:%v", st, et)
		startTime, err := strconv.ParseInt(st, 10, 64)
		if err != nil {
			log.Error("queryTimeReceipt timestamp err")
			return c.JSON(http.StatusInternalServerError, errors.New("queryTimeReceipt startTime timestamp err"))
		}
		endTime, err := strconv.ParseInt(et, 10, 64)

		res, err := rdb.ZRangeByScore(ctx, "TransactionSet", &redis.ZRangeBy{
			Min: strconv.FormatInt(startTime, 10),
			Max: strconv.FormatInt(endTime, 10),
		}).Result()

		if err != nil {
			log.Error("GET error: ", err)
			return c.JSON(http.StatusOK, res)
		}
		log.Info("查询交易数据成功: ")
		return c.JSON(http.StatusOK, res)
	})

	//接收到请求后根据数据类型进行分类，然后到指定路径查询账本内容返回，
	e.POST("/queryBlockInfos", func(c echo.Context) error {

		log.Info("按时间戳查询存证数据1")
		ledger := c.FormValue("blockType")
		c.FormValue("StartTime")
		start, err := time.ParseInLocation("2006-01-02 15:04:05", c.FormValue("StartTime"), LOC)
		if err != nil {
			log.Error("GET error: ", err)
			return c.JSON(http.StatusOK, nil)
		}
		time, index := backend.GetIndexMinInt(fmt.Sprint(start))
		log.Info(time, "+", index)
		//}

		var (
			fileName = "E:\\Go_WorkSpace\\hraft1102\\scope\\" + time + "\\" + ledger + "\\MINUTE" + "\\"
		)
		blockHeader := []backend.BlockHeader{}
		for i := 0; i < 20 && index-i > 0; i++ {
			tmp := ReadTxMinFiletoTenmin(fileName + strconv.Itoa(index-i))
			if tmp.KeyId == "" {
				continue
			}
			blockHeader = append(blockHeader, tmp)
		}
		//res, _ := json.Marshal(blockHeader)

		log.Info("查询区块头数据成功: ", blockHeader)
		return c.JSON(http.StatusOK, blockHeader)
	})
	return e
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
func ReadTxMinFiletoTenmin(fileName string) backend.BlockHeader {

	log.Info("区块文件名：", fileName)
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Info("区块文件不存在")
		return backend.BlockHeader{}
	} else {
		log.Info("文件存在，开始查询")
		jsonfile, _ := os.Open(fileName)
		defer jsonfile.Close()
		//读取文件
		fileContent, err := ioutil.ReadAll(jsonfile)
		if err != nil {
			log.Error("Read file err =: ", err)
		}
		//fmt.Println(string(fileContent))
		var minBlock backend.MinuteDataBlock
		//var user User
		if err := json.Unmarshal([]byte(fileContent), &minBlock); err != nil {
			log.Error("反解析 file error =: ", err)
		}

		//fmt.Println(minBlock.Header)
		//fmt.Println(len(minBlock.DataReceipts))
		return minBlock.Header
	}
}
