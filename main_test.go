package main

import (
	"cache/backend"
	"cache/iot_server"
	_"net"
	"testing"
	"time"
	"context"
	"cache/util/config"
	"fmt"
	"strconv"
	"encoding/json"
)

func TestReadTxMinFiletoTenmin(t *testing.T) {
	//var yearMonthDay string
	//yearMonthDay = time.Now().Format("2006-01-02")

	backend.ReadTxMinFiletoTenmin("2022-07-10","video","620")
}
func TestConvert(t *testing.T){
	tmp := backend.ReadTxMinFiletoTenmin("2022-07-10","video","606")
	end := backend.Convert(tmp.Header, tmp.DataReceipts[0])
	fmt.Println(*end)
	res, _ := json.Marshal(*end)
	fmt.Println(string(res))
}

func TestGetIndexMinInt(t *testing.T) {
	a, x := backend.GetIndexMinInt("2022-07-12 10:06:30.896002")
	fmt.Println(a, " ", x)
}
func TestGetValue(t *testing.T) {
	ctx := context.Background()
	config := config.Initialize()
	//交易信息打包推送到hraft节点
	//go backend.StartGrpcPort(":" + "8880")
	redisClient, err := backend.NewRedisBackend(ctx, &config)
	if err != nil {
		fmt.Println("后台程序运行失败: ", err)
		panic(err)
	}
	fmt.Println(time.Now().Unix())
//1657618674
	backend.GetValue(ctx, redisClient, time.Now().Unix())
}
func TestReadTxMinFiletoTenmin2(t *testing.T) {
	//var yearMonthDay string
	//yearMonthDay = time.Now().Format("2006-01-02")

	blockHeader := []backend.BlockHeader{}
	index:=1300
	for i:=0; i<20 && index-i>0 ; i++{
	//D:\\Go\\src\\hraft1102" + "/scope/" + time + "/" + ledger + "/MINUTE" + "/"
		tmp := iot_server.ReadTxMinFiletoTenmin("D:\\Go\\src\\hraft1102\\scope\\2022-07-21\\video\\MINUTE"+"/"+strconv.Itoa(index-i))
		fmt.Println(tmp)
		if tmp.KeyId == ""{
			continue
		}
		blockHeader = append(blockHeader,tmp)
	}
	//fmt.Println(blockHeader)
	res, _ := json.Marshal(blockHeader)
	fmt.Println(string(res))
	//res, _ := json.Marshal(blockHeader)

}
