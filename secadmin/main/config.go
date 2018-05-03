package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"fmt"
	"strings"
)

type MysqlConfig struct {
	UserName string
	Passwd   string
	Port     string
	Database string
	Host     string
}
type EtcdConf struct {
	Addr          string
	EtcdKeyPrefix string
	ProductKey    string
	Timeout       int
}

//小写得变量   在json解析中解析不出来得
type Config struct {
	mysqlConfig MysqlConfig
	etcdConf    EtcdConf
}

var (
	AppConf Config
)

func initConfig() (err error) {

	userName := beego.AppConfig.String("mysql_user_name")
	if len(userName) == 0 {
		logs.Error("mysql_user_name  is nil")
		return
	}
	passWd := beego.AppConfig.String("mysql_passwd")
	if len(passWd) == 0 {
		logs.Error("mysql_passwd is nil")
		return
	}
	host := beego.AppConfig.String("mysql_host")
	if len(host) == 0 {
		logs.Error("mysql_host is nil")
		return
	}
	dbName := beego.AppConfig.String("mysql_database")
	if len(dbName) == 0 {
		logs.Error("mysql_database is nil")
		return
	}
	port := beego.AppConfig.String("mysql_port")
	if len(port) == 0 {
		logs.Error("mysql_port is  nil")
		return
	}
	AppConf = Config{}

	AppConf.mysqlConfig.UserName = userName
	AppConf.mysqlConfig.Host = host
	AppConf.mysqlConfig.Database = dbName
	AppConf.mysqlConfig.Port = port
	AppConf.mysqlConfig.Passwd = passWd

	etcdAddr := beego.AppConfig.String("etcd_addr")
	if len(etcdAddr) == 0 {
		logs.Error("etcd_addr is nil")
		return
	}
	etcdKeyPrefix := beego.AppConfig.String("etcd_sec_key_prefix")
	if len(etcdKeyPrefix) == 0 {
		logs.Error("etcd_sec_key_prefix is  nil")
		return
	}

	etcdProductKey := beego.AppConfig.String("etcd_product_key")
	if len(etcdProductKey) == 0 {
		logs.Error("etcd_product_key is nil")
		return
	}
	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	if err != nil {
		logs.Error("etcd_timeout is failed ,err:%v", err)
		return
	}
	if !strings.HasSuffix(etcdKeyPrefix, "/") {
		etcdKeyPrefix = etcdKeyPrefix + "/"
	}
	AppConf.etcdConf.Addr = etcdAddr
	AppConf.etcdConf.EtcdKeyPrefix = etcdKeyPrefix
	AppConf.etcdConf.ProductKey = fmt.Sprintf("%s%s",AppConf.etcdConf.EtcdKeyPrefix, etcdProductKey)
	AppConf.etcdConf.Timeout = etcdTimeout

	logs.Debug("init config succ")

	logs.Debug(AppConf)
	return
}
