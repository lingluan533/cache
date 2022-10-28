package backend

import (
	rpc "cache/backend/rpc/proto"
	"context"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"testing"
)
func GetQueryServiceClient() rpc.QueryServiceClient {
	serviceAddress := "10.128.209.25:8888"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	Client := rpc.NewQueryServiceClient(conn)
	return Client
}
func TestBlockService_QueryBlock(t *testing.T) {
	client:= GetQueryServiceClient()
	var req= &rpc.RequestBlock{
		LedgerType:"service_access",
		BlockChainType:"minute",
		Height:1,
		KeyId: "0001",
	}
	resp,err:=client.QueryBlock(context.Background(),req)
	if err!=nil{
		t.Fatal(err)
	}
	var block=&rpc.MinuteTxBlock{}
	proto.Unmarshal(resp.Block,block)
	t.Logf("%+v\n",*block)
}
func TestBlockService_QueryBlockBatch(t *testing.T) {
	client:= GetQueryServiceClient()
	var req= &rpc.RequestDataBatch{
		LedgerType:"service_access",
		BlockChainType:"minute",
		DataType:"transaction",
		Num: 10,
	}
	results,err:=client.QueryBlockBatch(context.Background(),req)
	if err!=nil{
		t.Fatal(err)
	}
	for _,data:=range results.Data {
		t.Logf("%+v\n", string(data))
	}
}
func TestBlockService_QueryTx(t *testing.T) {
	client:= GetQueryServiceClient()
	var req= &rpc.RequestTx{
		LedgerType:"service_access",
		BlockChainType:"minute",
		Height:1,
		TransactionId:"0001",
	}
	result,err:=client.QueryTx(context.Background(),req)
	if err!=nil{
		t.Fatal(err)
	}
	t.Logf("%+v\n",*result)
}
func TestBlockService_QueryDataReceipt(t *testing.T) {
	client:= GetQueryServiceClient()
	var req= &rpc.RequestDataReceipt{
		LedgerType:"video",
		BlockChainType:"minute",
		Height:1,
		KeyId:"123",
	}
	result,err:=client.QueryDataReceipt(context.Background(),req)
	if err!=nil{
		t.Fatal(err)
	}
	t.Logf("%+v\n",*result)
}
func TestBlockService_QueryGenesisBlock(t *testing.T) {
	client:= GetQueryServiceClient()
	var req= &rpc.RequestGenesisBlock{
		LedgerType:"service_access",
		BlockChainType:"minute",
		Hash:"0x7823123",
	}
	block,err:=client.QueryGenesisBlock(context.Background(),req)
	if err!=nil{
		t.Fatal(err)
	}
	t.Logf("%+v\n",*block)
}

func TestReadTxMinFiletoTenmin(t *testing.T) {
	
}
