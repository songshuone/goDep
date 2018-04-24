package main

import (
	"os"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	etcd "github.com/coreos/etcd/clientv3"
	"time"
	"context"
	"goDep/skcproxy/service"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

var (
	redisPool  *redis.Pool
	etcdClient *etcd.Client
)

func initSec() (err error) {
	err = initLogger()
	if os.IsExist(err) {
		return
	}
	logs.Info("init logger success!")
	err = initEtcd()
	if os.IsExist(err) {
		return
	}
	logs.Info("init etcd success!")
	//err = initRedis()
	//if os.IsExist(err) {
	//	return
	//}
	//logs.Info("init redis success!")

	service.IntentService(service.Conf)

	err = loadSecConf()
	if os.IsExist(err) {
		return
	}
	logs.Info("init secConf success!")

	initSecProductWatcher()
	return
}

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initLogger() (err error) {
	config := make(map[string]interface{})

	config["filename"] = service.Conf.LogPath
	config["level"] = convertLogLevel(service.Conf.LogLevel)

	configStr, err := json.Marshal(config)

	if os.IsExist(err) {
		fmt.Println("marsha1 failed err:", err)
	}
	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}
func initEtcd() (err error) {
	etcdClient, err = etcd.New(etcd.Config{Endpoints: []string{service.Conf.EtcdAddr}, DialTimeout: time.Duration(service.Conf.EtcdTimeout) * time.Second,})
	if os.IsExist(err) {
		logs.Error("connect etcd failed err:", err)
		return
	}

	return
}
//func initRedis() (err error) {
//
//	redisPool = &redis.Pool{MaxIdle: service.Conf.RedisMaxIdle, MaxActive: service.Conf.RedisMaxActive, IdleTimeout: time.Duration(service.Conf.RedisIdleTimeout) * time.Second, Dial: func() (redis.Conn, error) {
//		return redis.Dial("tcp", service.Conf.RedisAddr)
//	}}
//	conn := redisPool.Get()
//
//	defer conn.Close()
//
//	_, err = conn.Do("ping")
//
//	if os.IsExist(err) {
//		logs.Error("ping redis failed", err)
//		return
//	}
//	return
//}

func loadSecConf() (err error) {

	resp, err := etcdClient.Get(context.Background(), service.Conf.EtcdProductKey)
	if os.IsExist(err) {
		logs.Error(fmt.Sprintf("get [%v] from etcd failed ,err%v", service.Conf.EtcdProductKey, err))
		return
	}

	var secProductInfo [] service.SecProductInfoConf

	for k, v := range resp.Kvs {
		logs.Debug("key[%v] values[%v]", k, v)
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			logs.Error("Unmarshal sec product info failed ,err:", err)
			return
		}
		logs.Debug("sec info conf is [%v]", secProductInfo)
	}

	updateSecProductInfo(secProductInfo)
	return
}

//更新秒杀的数据
func updateSecProductInfo(secProductInfo [] service.SecProductInfoConf) {
	var tmp map[int]*service.SecProductInfoConf = make(map[int]*service.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {

		//每次v都是指向同一个地址 但是每次地址的值不一样  会保留最后一个v的值
		//tmp[v.ProductId] = &v          会出现所有的tmp的值的地址指向同一个 （v最后出现的值）

		//修改
		temp := v
		tmp[v.ProductId] = &temp

	}
	service.Conf.RWSecProductLock.Lock()
	service.Conf.SecProductInfoMap = tmp
	service.Conf.RWSecProductLock.Unlock()
}
func initSecProductWatcher() {
	go watchSecProductKey(service.Conf.EtcdProductKey)
}

//监听etcd的写入的数据
func watchSecProductKey(etcdKey string) {

	//cli, err := etcd.New(etcd.Config{Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}, DialTimeout: time.Duration(model.Conf.EtcdTimeout) * time.Second})
	//if err != nil {
	//	logs.Error("connect ectd failed ,err[%v]", err.Error())
	//	return
	//}

	logs.Debug("begin watch key :[%s]", etcdKey)

	for {
		rch := etcdClient.Watch(context.Background(), etcdKey)

		var getConSucc = true
		var secProduct []service.SecProductInfoConf
		for wresp := range rch {

			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%v] 's config deleted", etcdKey)
					continue
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == etcdKey && ev.Kv.Value != nil {

					err := json.Unmarshal(ev.Kv.Value, &secProduct)
					logs.Warn(string(ev.Kv.Value), "==================", secProduct)
					if err != nil {
						logs.Error("key [%s],Unmarshal ,err:%v", etcdKey, err)
						getConSucc = false
						continue
					}
				}
				logs.Debug("get config form etcd %s %q:%q", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConSucc {
				logs.Debug("get config from etcd succ,%v", secProduct)
				updateSecProductInfo(secProduct)
			}
		}
	}

}
