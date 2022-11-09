package main

import (
	"bytes"
	"cache/iot_server"
	"encoding/json"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

//在线服务列表
var OnlineServers []*consulapi.ServiceEntry

func main() {
	//定时更新在线服务列表
	//ticker1 := time.NewTicker(time.Second * 5)
	//go func() {
	//	for {
	//		<-ticker1.C
	//		OnlineServers = consul_service.QueryOnlineEdgeNodes()
	//		ticker1.Reset(time.Minute * 1)
	//	}
	//}()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, iot_server.NewResult("hello, transfer server!"))
	})
	e.POST("/postIoTData", func(c echo.Context) error {
		var transactions iot_server.DataTransactions
		if err := c.Bind(&transactions); err != nil {
			fmt.Println("marshal error: ", err)
			return c.JSON(http.StatusOK, iot_server.NewResult(err.Error(), nil))
		}
		fmt.Println(transactions)

		if len(OnlineServers) == 0 || OnlineServers == nil {
			return c.JSON(http.StatusInternalServerError, "No Avaliable EdgeNode!")
		}
		service := OnlineServers[0].Service
		byte, err := json.Marshal(transactions)
		resp, err := http.Post("http://"+service.Address+":"+strconv.Itoa(service.Port)+"/storeTransaction", "application/json", bytes.NewBuffer(byte))

		//resp, err := http.PostForm("http://"+service.Address+":"+strconv.Itoa(service.Port)+"/storeTransaction", url.Values{"createTimestamp": {transactions.CreateTimestamp}, "entityId": {transactions.EntityId}, "turnOver": {strconv.FormatFloat(transactions.TurnOver, 'E', 2, 64)}, "txRecNum": {strconv.Itoa(int(transactions.TxRecNum))}, "transactions": {""}})

		if err != nil {
			fmt.Printf("Error on request: %v\n", err)
			return c.JSON(http.StatusInternalServerError, "Unmarshalerr error")
		}
		defer resp.Body.Close()
		return c.String(http.StatusOK, "Hello, Scope!")
	})
	e.Start(":9001")
}
