package service

import (
	"github.com/astaxie/beego/logs"
	"fmt"
	"time"
	"crypto/md5"
)



func NewSecRequest() (secRequest *SecRequest) {
	secRequest = &SecRequest{
		ResultChan: make(chan *SecResult),
	}
	return
}

//用户检查
func userCheck(secReuest *SecRequest) (err error) {

	found := false
	//白名单校验
	for _, v := range Conf.ReferWhiteList {
		if v == secReuest.ClientRefence {
			found = true
			break
		}
	}
	if found {
		err = fmt.Errorf("invalid request")
		logs.Warn("userId[%d] is reject by refer,req[%v]", secReuest.UserId, secReuest)
		return
	}

	//cookie 校验
	authData := fmt.Sprintf("%d:%s", secReuest.UserId, secKillConf.CookieSecretKey)

	//%x  16进制  %x
	// 十六进制表示，字母形式为小写 a-f         Printf("%x", 13)             d
	//%X      十六进制表示，字母形式为大写 A-F
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))
    logs.Debug("authSign:%s ,secReuest.UserAuthSign:%s",authSign,secReuest.UserAuthSign)
	if authSign != secReuest.UserAuthSign {
		err = fmt.Errorf("invalid cookie:userId")
		return
	}
	return
}

//开始秒杀
func SecKill(secReuest *SecRequest) (data map[string]interface{}, code int, err error) {
	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.RUnlock()

	err = userCheck(secReuest)
	if err != nil {
		code = ErrUserCheckAuthFailed
		logs.Warn("userId[%v] invalid ,check failed request[%v]", secReuest.UserId, secReuest)
		return
	}
	//限制恶意用户
	err = antiSpam(secReuest)
	if err != nil {
		code = ErrUserServiceBusy
		logs.Warn("userId[%v] invalid ,check failed request[%v]", secReuest.UserId, secReuest)
		return
	}

	data, code, err = SecInfoById(secReuest.ProductId)

	//获取商品失败
	if err != nil {
		logs.Warn("userId[%d] secInfoBy Id failed, req[%v]", secReuest.UserId, secReuest)
		return
	}
	if code != 0 {
		logs.Warn("userId[%d] secInfoBy Id failed, req[%v]", secReuest.UserId, secReuest)
		return
	}

	userKey := fmt.Sprintf("%d_%d", secReuest.UserId, secReuest.ProductId)

	secKillConf.UserConnMap[userKey] = secReuest.ResultChan
	logs.Debug("userKey:=========",userKey)
	secKillConf.SecReqChan <- secReuest


	//定时10s
	ticker := time.NewTicker(time.Second * 10)

	defer func() {
		ticker.Stop()
		secKillConf.UserConnMapLock.Lock()
		//删除map中的值
		delete(secKillConf.UserConnMap, userKey)
		secKillConf.UserConnMapLock.Unlock()
	}()
	select {
	case <-ticker.C:
		//连接超时
		code = ErrProcessTimeout
		err = fmt.Errorf("request timeout")
		return
	case <-secReuest.CloseNotify:
		//客户端关闭连接
		code = ErrClientClosed
		err = fmt.Errorf("client already colsed")
		return
	case result := <-secReuest.ResultChan:
		code = result.Code
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_Id"] = result.UserId
		return
	}

	return
}



//获取秒杀的商品
func SecInfo(productId int) (data []map[string]interface{}, code int, err error) {

	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.RUnlock()
	v, ok := secKillConf.SecProductInfoMap[productId]
	data = make([]map[string]interface{}, 0)
	if !ok {
		for _, v := range secKillConf.SecProductInfoMap {
			item, _, err := SecInfoById(v.ProductId)
			if err != nil {
				logs.Error("productId[%v] is nil", v.ProductId)
				continue
			}
			logs.Debug("productId[%v] map[%v]", v.ProductId, v)
			data = append(data, item)
		}
		return
	}
	item, code, err := SecInfoById(v.ProductId)
	if err != nil {
		return
	}
	data = append(data, item)
	return
}

func SecInfoById(productId int) (data map[string]interface{}, code int, err error) {

	secKillConf.RWSecProductLock.RLock()
	defer secKillConf.RWSecProductLock.RUnlock()
	v, ok := secKillConf.SecProductInfoMap[productId]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found product_id:%d", productId)
		return
	}
	start := false
	end := false
	state := "success"

	now := time.Now().Unix()

	if now-v.StartTime < 0 {
		//还没有开始
		state = "sec kill is not start"
		code = ErrActiveNotStart
	}

	if now-v.StartTime > 0 {
		//已经 开始
		start = true
	}

	if now-v.EndTime > 0 {
		//已经结束
		start = false
		end = true
		state = "sec kill is already end"
		code = ErrActiveAlreadyEnd
	}

	if v.Status == ProductStatusSaleOut || v.Status == ProductStatusForceSaleOut {
		//抢完或强制关闭
		start = false
		end = true
		state = "sec kill is sale out"
		code = ErrActiveSaleOut
	}

	data = make(map[string]interface{})
	data["product_id"] = productId
	data["start"] = start
	data["end"] = end
	data["status"] = state
	//data["total"] = v.Total
	//data["left"] = v.Left
	return
}
