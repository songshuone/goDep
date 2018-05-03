package main

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	etcd_client "github.com/coreos/etcd/clientv3"
	"time"
	"goDep/secadmin/model"
	"github.com/jinzhu/gorm"
)

var (
	EtcdClient *etcd_client.Client
	Db         *gorm.DB
)

func initEtcd() (err error) {
	client, err := etcd_client.New(etcd_client.Config{
		Endpoints: []string{AppConf.etcdConf.Addr}, DialTimeout: time.Duration(AppConf.etcdConf.Timeout) * time.Second,
	})
	if err != nil {
		logs.Error("init etcd failed ,err:%v", err)
		return
	}
	EtcdClient = client
	logs.Debug("init etcd succ")
	return
}

func initDb() (err error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local", AppConf.mysqlConfig.UserName, AppConf.mysqlConfig.Passwd, AppConf.mysqlConfig.Host, AppConf.mysqlConfig.Port, AppConf.mysqlConfig.Database)
	logs.Debug(dns)
	database, err := gorm.Open("mysql", dns)
	if err != nil {
		logs.Error("open mysql failed err:%v", err)
		return
	}
	Db = database
	logs.Debug("connect to mysql succ")
	return
}

func initAll() (err error) {

	err = initConfig()
	if err != nil {
		logs.Error("init config failed err:%v", err)
		return
	}

	err = initDb()
	if err != nil {
		logs.Error("init db failed err:%v", err)
		return
	}

	err = initEtcd()

	if err != nil {
		logs.Error("init etcd failed, err:%s", err)
		return
	}

	err = model.Init(Db, EtcdClient,  AppConf.etcdConf.ProductKey,AppConf.etcdConf.EtcdKeyPrefix)
	if err != nil {
		logs.Error("init  model failed, err:%s", err)
		return
	}
	return
}
