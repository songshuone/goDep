package service

import (
	"github.com/astaxie/beego/logs"
	"context"
	"time"
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

func initSecProductWatch(con *SecLayerConf) {
	go watchSecInfo(con)
}

//监听 etcd
func watchSecInfo(con *SecLayerConf) {
	etcdKey := con.EtcdConfig.EtcdProductKey
	for {
		var secProductInfo []SecProductInfoConf
		getConfSucc := true
		wc := secLayerContext.etcdClient.Watch(context.Background(), etcdKey)

		for wresp := range wc {
			for _, ev := range wresp.Events {
				//监听到删除
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", etcdKey)
					continue
				}
				//监听到写入
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == etcdKey {
					err := json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", secProductInfo)
				updateSecProductInfo(con, secProductInfo)
			}
		}
	}
}

func updateSecProductInfo(con *SecLayerConf, secProductInfo []SecProductInfoConf) {
	var tmp map[int]*SecProductInfoConf = make(map[int]*SecProductInfoConf, 1024)

	for _, v := range secProductInfo {
		productInfo := v
		productInfo.secLimit = &SecLimit{}
		tmp[productInfo.ProductId] = &productInfo
	}
	secLayerContext.RWSecProductLock.Lock()
	con.SecProductInfoMap = tmp
	secLayerContext.secLayerConf = con
	secLayerContext.secLayerConf.SecProductInfoMap = con.SecProductInfoMap
	secLayerContext.RWSecProductLock.Unlock()
	logs.Debug("secLayerContext.secLayerConf.SecProductInfoMap :", secLayerContext.secLayerConf.SecProductInfoMap)
}

func loadProductFromEtcd(con *SecLayerConf) (err error) {
	logs.Debug("start from etcd read data")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := secLayerContext.etcdClient.Get(ctx, con.EtcdConfig.EtcdProductKey)
	if err != nil {
		logs.Error("get[%s] from etcd failed,err:%v", con.EtcdConfig.EtcdProductKey, err)
		return
	}
	logs.Debug("get [%s] from ectd success,resp:%v", con.EtcdConfig.EtcdProductKey, resp)

	var secProductInfo []SecProductInfoConf

	for k, v := range resp.Kvs {
		logs.Debug("key[%v],value[%v]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("json Unmarshal failed,err:%v", err)
			return
		}
		logs.Debug("sec info conf is [%v]", secProductInfo)
	}

	updateSecProductInfo(con, secProductInfo)

	logs.Debug("update  sec product info success")

	initSecProductWatch(con)

	logs.Debug("init etcd watch success!!!")
	return
}
