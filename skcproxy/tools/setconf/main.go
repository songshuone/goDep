package main

import (
	etcd "github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"encoding/json"
	"context"
	"goDep/skcproxy/service"
)

const (
	EtcdKey = "/backend/secskill/product"
)

func setLogConfToEtcd() {

	cli, err := etcd.New(etcd.Config{Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}, DialTimeout: 5 * time.Second,})

	if err != nil {
		fmt.Println("connect failed, err :", err)
		panic(err)
	}

	fmt.Println("connect etcd success")

	defer cli.Close()

	var secProductInfo []service.SecProductInfoConf

	secProductInfo = append(secProductInfo,
		service.SecProductInfoConf{
			ProductId:         3000,
			StartTime:         1524534841,
			EndTime:           1524634841,
			Status:            0,
			Total:             1000,
			Left:              1000,
			OnePersonBuyLimit: 1,
			BuyRate:           0.5,
			SoldMaxLimit:      100,}, )

	now := time.Now().Unix()
	secProductInfo = append(
		secProductInfo,
		service.SecProductInfoConf{
			ProductId:         2000,
			StartTime:         now,
			EndTime:           now + 1505012400,
			Status:            0,
			Total:             2000,
			Left:              1000,
			OnePersonBuyLimit: 50,
			BuyRate:           0.8,
			SoldMaxLimit:      100,
		},
	)


	data, err := json.Marshal(secProductInfo)
	if err != nil {
		fmt.Println("json failed,", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	_, err = cli.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		fmt.Println("put failed,err", err)
		return
	}

	//ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//
	//resp, err := cli.Get(ctx, EtcdKey)
	//cancel()
	//if err != nil {
	//	fmt.Println("get failed,err", err)
	//	return
	//}
	//for _, ev := range resp.Kvs {
	//	fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	//}
}

func main() {
	setLogConfToEtcd()
}
