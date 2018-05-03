package model

import (
	etcd_client "github.com/coreos/etcd/clientv3"
	"github.com/jinzhu/gorm"
)

var (
	EtcdClient     *etcd_client.Client
	EtcdPrefix     string
	EtcdProductKey string
	Db             *gorm.DB
)

func Init(db *gorm.DB, etcdClient *etcd_client.Client, etcdProductKey string, etcdPrefix string) (err error) {
	EtcdProductKey = etcdProductKey
	EtcdClient = etcdClient
	EtcdPrefix = etcdPrefix
	Db = db
	err = CreateTable()
	return
}

func CreateTable() (err error) {
	err = Db.AutoMigrate(new(Product)).Error
	if err != nil {
		return err
	}
	err = Db.AutoMigrate(new(Activity)).Error
	if err != nil {
		return err
	}
	return
}
