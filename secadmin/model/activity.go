package model

import (
	"time"
	"fmt"
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

const (
	ActivityStatusNormal  = 0
	ActivityStatusDisable = 1
	ActivityStatusExpire  = 2
)

type ActivityModel struct {
}

type Activity struct {
	ActivityId   int `gorm:"primary_key"`
	ActivityName string
	ProductId    int
	StartTime    int64
	EndTime      int64
	Total        int
	Status       int

	StartTimeStr string
	EndTimeStr   string
	StatusStr    string
	Speed        float64 `gorm:"column:sec_speed"`
	BuyLimit     int     `gorm:"column:buy_limit"`
}

type SecProductInfoConf struct {
	ProductId         int     `json:"product_id"`
	StartTime         int64   `json:"start_time"`
	EndTime           int64   `json:"end_time"`
	Status            int     `json:"status"`
	Total             int     `json:"total"`
	Left              int     `json:"left"`
	OnePersonBuyLimit int     `json:"one_person_buy_limit"`
	BuyRate           float64 `json:"buy_rate"`
	//每秒最多能卖多少个
	SoldMaxLimit int `json:"sold_max_limit"`
}

func NewActivityModel() (*ActivityModel) {
	return &ActivityModel{}
}

func CheckProductIsExits(product int) bool {
	if err := Db.Find(new(Product), product).Error; err != nil {
		return false
	} else {
		return true
	}

}

func loadProductFromEtcd() (secProductInfos []SecProductInfoConf, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	gr, err := EtcdClient.Get(ctx, EtcdProductKey)
	if err != nil {
		logs.Warn("get product list failed from etcd,err:%v", err)
		return
	}
	secProductInfos = make([]SecProductInfoConf, 0)
	for _, v := range gr.Kvs {
		err = json.Unmarshal(v.Value, &secProductInfos)
		if err != nil {
			logs.Error("json unmarshl failed,err:%v", err)
			return
		}
		logs.Debug("get product form etcd [%s]info:%v", EtcdProductKey,secProductInfos)
		return
	}

	return
}

func SendEtcd(activity *Activity) (err error) {

	secProductInfos, err := loadProductFromEtcd()
	if err != nil {
		return
	}
	secProductInfo := SecProductInfoConf{}
	secProductInfo.ProductId = activity.ProductId
	secProductInfo.EndTime = activity.EndTime
	secProductInfo.OnePersonBuyLimit = activity.BuyLimit
	secProductInfo.ProductId = activity.ProductId
	secProductInfo.BuyRate = activity.Speed
	secProductInfo.SoldMaxLimit = activity.Total / 10
	secProductInfo.StartTime = activity.StartTime
	secProductInfo.Status = activity.Status
	secProductInfo.Total = activity.Total

	secProductInfos = append(secProductInfos, secProductInfo)
	data, err := json.Marshal(secProductInfos)
	if err != nil {
		err = fmt.Errorf("josn marshal failed,err:%v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = EtcdClient.Put(ctx, EtcdProductKey, string(data))
	if err != nil {
		return
	}
	logs.Debug("send data to etcd success[%s]:[%s]",EtcdProductKey, string(data))
	return
}
func (a *ActivityModel) CreateActivity(activity *Activity) (err error) {

	if CheckProductIsExits(activity.ProductId) {
		db := Db.Begin()
		db.Create(activity)
		err = SendEtcd(activity)
		if err != nil {
			db.Rollback()
			return
		}
		db.Commit()

	} else {
		err = fmt.Errorf("该商品Id %d不存在", activity.ProductId)
	}

	return
}

func (p *ActivityModel) GetActivityList() (activity []*Activity, err error) {
	err = Db.Find(&activity).Error

	for _, v := range activity {
		t := time.Unix(v.StartTime, 0)
		v.StartTimeStr = t.Format("2006-01-02 15:04:05")

		t = time.Unix(v.EndTime, 0)
		v.EndTimeStr = t.Format("2006-01-01 15:04:05")

		now := time.Now().Unix()

		if now > v.EndTime {
			v.StatusStr = "已结束"
			continue
		}

		if v.Status == ActivityStatusNormal {
			v.StatusStr = "正常"
		} else if v.Status == ActivityStatusDisable {
			v.StatusStr = "已禁用"
		}

	}

	return

}
