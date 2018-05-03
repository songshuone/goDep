package service

//条件逻辑处理时   先处理错误的情况    然后在处理正确的数据
import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"os"
	"github.com/garyburd/redigo/redis"
	"fmt"
)

//将秒杀到的结果存入redis中   交给第三方服务确认
func WriteHandle() {

	for {
		req := <-secKillConf.SecReqChan
		data, err := json.Marshal(*req)
		if err != nil {
			logs.Error("json Marshal failed err:%v,req[%v]====", err, req)
			continue
		}
		conn := secKillConf.proxy2LayerRedisPool.Get()
		//从队列的左边入队一个或多个元素
		_, err = conn.Do("LPUSH", "redis_proxy2layer_queue_name", string(data))
		if os.IsExist(err) {
			logs.Error("lpush failed, err:%v ,req[%v]", err, req)
			conn.Close()
			continue
		}
		conn.Close()
	}

}

//三方服务确认 之后   获取被秒杀到的用户
func ReadHandle() {
	for {
		conn := secKillConf.proxy2LayerRedisPool.Get()
		//从队列的右边取一个元素
		reply, err := conn.Do("RPOP", "redis_layer2proxy_queue_name")
		data, err := redis.String(reply, err)

		if os.IsExist(err) {
			logs.Error("rpop failed, err:%v", err)
			conn.Close()
			continue
		}
		var result SecResult
		err = json.Unmarshal([]byte(data), &result)

		//logs.Debug("result", result)
		if os.IsExist(err) {
			logs.Error("json Unmarshal failed,err:%v", err)
			conn.Close()
			continue
		}

		secKillConf.UserConnMapLock.Lock()
		userKey := fmt.Sprintf("%d_%d", result.UserId, result.ProductId)
		resultChan, ok := secKillConf.UserConnMap[userKey]
		secKillConf.UserConnMapLock.Unlock()
		if !ok {
			conn.Close()
			//logs.Warn("user not found:%v", userKey)
			continue
		}

		resultChan <- &result
		conn.Close()

	}
}
