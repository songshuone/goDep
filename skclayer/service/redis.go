package service

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"math/rand"
	"crypto/md5"
)

//读取代理层请求的数据
func HandleReader() {

	logs.Debug("read goroutine running")

	for {
		conn := secLayerContext.proxy2LayerRedisPool.Get()
		for {
			reply, err := conn.Do("blpop", secLayerContext.secLayerConf.Proxy2LayerRedis.RedisQueueName, 0)
			data, err := redis.Strings(reply, err)

			if err != nil {
				logs.Error("pop from queue failed ,err:%v", err)
				break
			}
			logs.Debug("pop from queue , data:%s,length=%d", data, len(data))

			var req SecRequest

			err = json.Unmarshal([]byte(data[1]), &req)

			if err != nil {
				logs.Error("json Unmarshal failed err:%v", err)
				continue
			}
			logs.Debug("pop from queue , req:%v", req)

			now := time.Now().Unix()

			if now-req.AccessTime.Unix() >= int64(secLayerContext.secLayerConf.MaxRequestWaitTimeout) {
				logs.Warn("req[%v] is expire", req)
				continue
			}
			timer := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.secLayerConf.SendToHandleChanTimeout))
			select {
			case secLayerContext.Read2HandleChan <- &req:

			case <-timer.C:
				logs.Warn("send to handle chan timeout,err:%v", req)
				break
			}
		}
		conn.Close()
	}

}

func sendToRedis(res *SecResponse) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		logs.Error("json Marshal failed ,err:%v", err)
		return
	}

	logs.Debug("secLayerContext.Handle2WriteChan", string(data))
	conn := secLayerContext.layer2ProxyRedisPool.Get()
	_, err = conn.Do("rpush", secLayerContext.secLayerConf.Layer2ProxyRedis.RedisQueueName, string(data),0)
	if err != nil {
		logs.Warn("rpush redis failed ,err:%v", err)
		conn.Close()
		return
	}
	conn.Close()
	return
}

//把处理抢购完成的接口send redis 中
func HandleWrite() {
	logs.Debug("handle write running")
	for {
		res := <-secLayerContext.Handle2WriteChan
		logs.Debug("secLayerContext.Handle2WriteChan", res)
		err := sendToRedis(res)
		if err != nil {
			logs.Error("send redis failed,err:%v", err)
			continue
		}
	}

}

//从代理层获取的数据进行抢购逻辑处理
func HandleUser() {
	logs.Debug("handle user running")
	for req := range secLayerContext.Read2HandleChan {
		logs.Debug("begin process request:%v", req)
		res, err := HandlerSecKill(req)
		if err != nil {
			logs.Warn("process request %v failed,err:%v", err)
			res = &SecResponse{
				Code: ErrServiceBusy,
			}
		}
		timer := time.NewTicker(time.Millisecond * time.Duration(secLayerContext.secLayerConf.SendToWriteChanTimeout))

		select {
		case <-timer.C:
			break
			//把抢购完成后的结果放入chan 管道中
		case secLayerContext.Handle2WriteChan <- res:
		}

	}
}

func HandlerSecKill(req *SecRequest) (res *SecResponse, err error) {

	secLayerContext.RWSecProductLock.Lock()
	defer secLayerContext.RWSecProductLock.Unlock()

	res = &SecResponse{}

	product, ok := secLayerContext.secLayerConf.SecProductInfoMap[req.ProductId]
	if !ok {
		logs.Error("not found product:%v", req.ProductId)
		res.Code = ErrNotFoundProduct
		return
	}

	//已经抢完
	if product.Status == ProductStatusSoldout {
		res.Code = ErrSoldout
		return
	}

	now := time.Now().Unix()
	//检查每秒最多能卖多少个的数量
	alreadySoldCount := product.secLimit.Check(now)

	if alreadySoldCount >= product.SoldMaxLimit { //如果大于 设置的  则重试
		res.Code = ErrRetry
		fmt.Println("ErrRetry===144")
		return
	}

	secLayerContext.HistoryMapLock.Lock()
	userHistory, ok := secLayerContext.HistoryMap[req.ProductId]
	if !ok {
		userHistory = &UserBuyHistory{history: make(map[int]int, 16)}
		secLayerContext.HistoryMap[req.ProductId] = userHistory
	}
	//获取一个抢的数量
	historyCount := userHistory.GetProductBuyCount(req.ProductId)
	secLayerContext.HistoryMapLock.Unlock()
	if historyCount >= product.OnePersonBuyLimit { //数量限制
		res.Code = ErrAlreadyBuy
		return
	}

	//获取当前商品已经被秒的数量
	curSoldCount := secLayerContext.productCountMgr.Count(req.ProductId)
	//已经被秒完
	if curSoldCount >= product.Total { //校验
		res.Code = ErrSoldout
		product.Status = ProductStatusSoldout
		return
	}

	//被秒的随机
	curRate := rand.Float64()
	fmt.Println(curRate,"============",product.BuyRate)
	if curRate > product.BuyRate {
		//没抢到  请重试
		res.Code = ErrRetry
		fmt.Println("ErrRetry 176")
		return
	}

	userHistory.Add(req.ProductId, 1)

	secLayerContext.productCountMgr.Add(req.ProductId, 1)

	res.Code = ErrSecKillSucc

	tokenData := fmt.Sprintf("userId=%d&productId=%d&timestamp=%d&security=%s", req.UserId, req.ProductId, now, secLayerContext.secLayerConf.TokenPasswd)

	res.Token = fmt.Sprintf("%x", md5.Sum([]byte(tokenData)))
	res.TokenTime = now
	res.ProductId=req.ProductId
	res.UserId=req.UserId
	fmt.Println("========================",res)
	return
}
func RunProcess() (err error) {
	for i := 0; i < secLayerContext.secLayerConf.ReadGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleReader()
	}

	for i := 0; i < secLayerContext.secLayerConf.ReadGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleUser()
	}

	for i := 0; i < secLayerContext.secLayerConf.ReadGoroutineNum; i++ {
		secLayerContext.waitGroup.Add(1)
		go HandleWrite()
	}

	secLayerContext.waitGroup.Wait()
	return
}

func initRedis(conf *SecLayerConf) (err error) {
	secLayerContext.proxy2LayerRedisPool, err = initRedisPool(conf.Proxy2LayerRedis)
	if err != nil {
		err = fmt.Errorf("init redis proxy2layer failed,err:%v", err)
		logs.Error(err)
		return
	}
	secLayerContext.layer2ProxyRedisPool, err = initRedisPool(conf.Layer2ProxyRedis)
	if err != nil {
		err = fmt.Errorf("init redis layer2proxy failed,err:%v", err)
		logs.Error(err)
		return
	}
	return
}

func initRedisPool(redisConf RedisConf) (redisPool *redis.Pool, err error) {
	redisPool = &redis.Pool{MaxIdle: redisConf.RedisMaxIdle, MaxActive: redisConf.RedisMaxActive, IdleTimeout: time.Duration(redisConf.RedisIdleTimeout) * time.Second, Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", redisConf.RedisAddr)
	}}

	conn := redisPool.Get()
	_, err = conn.Do("ping")
	defer conn.Close()
	if err != nil {
		logs.Error("ping redis failed", err)
		return
	}
	return
}
