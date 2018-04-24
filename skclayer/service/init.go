package service

import (
	"github.com/astaxie/beego/logs"
	etcd_client "github.com/coreos/etcd/clientv3"
	"time"
)

func initEtcd(config EtcdConfig) (err error) {
	cli, err := etcd_client.New(etcd_client.Config{Endpoints: []string{config.EtcdServiceAddr}, DialTimeout: time.Duration(config.EtcdTimeOut) * time.Second})
	if err != nil {
		logs.Error("connect etcd  failed ,err:%v", err)
		return
	}
	secLayerContext.etcdClient = cli
	logs.Debug("init ectd success!!")
	return
}

//初始化秒杀逻辑
func InitSecLayer(conf *SecLayerConf) (err error) {

	err = initRedis(conf)
	if err != nil {
		logs.Error("init redis failed ,err:%v", err)
		return
	}
	logs.Debug("init redis success!!")

	err = initEtcd(conf.EtcdConfig)

	if err != nil {
		logs.Error("init etcd failed,err:%v", err)
		return
	}

	err = loadProductFromEtcd(conf)
	if err != nil {
		logs.Error("load product from etcd failed,err:%v", err)
	}

	logs.Debug("init product success!!")

	secLayerContext.secLayerConf = conf
	secLayerContext.Read2HandleChan = make(chan *SecRequest, secLayerContext.secLayerConf.Read2handleChanSize)
	secLayerContext.Handle2WriteChan = make(chan *SecResponse, secLayerContext.secLayerConf.Handle2WriteChanSize)

	secLayerContext.HistoryMap = make(map[int]*UserBuyHistory, 1000000)
	secLayerContext.productCountMgr = NewProductCountMgr()

	logs.Debug("init all success!")

	return
}
