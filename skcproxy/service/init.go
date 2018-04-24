package service

import (
	"github.com/astaxie/beego/logs"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	secKillConf *SecSkillConf
)

//黑名单ip
func SyncIpBlackList() {

	var ipList []string

	lastTime := time.Now().Unix()
	for {
		conn := secKillConf.blackRedisPool.Get()
		//BLPOP  删除，并获取该列表中的第一个元素，或阻赛，知道有一个可用
		reply, err := conn.Do("BLPOP", "blackiplist", int(time.Second))

		ip, err := redis.Strings(reply, err)
		if err != nil {
			conn.Close()
			continue
		}

		curTime := time.Now().Unix()
		ipList = append(ipList, ip...)

		if len(ipList) > 2 || curTime-lastTime > 5 {
			secKillConf.RWBlackLock.Lock()

			for _, v := range ipList {
				secKillConf.IpBlackMap[v] = true
			}

			secKillConf.RWBlackLock.Unlock()
			lastTime = curTime

			logs.Info("sync ip list form redis succ,ip[%v],secKillConf[%v]", ipList,secKillConf)

		}

		conn.Close()

	}

}

//黑名单id
func SyncIdBlackList() {

	for {

		conn := secKillConf.blackRedisPool.Get()
		//BLPOP  删除，并获取该列表中的第一个元素，或阻赛，知道有一个可用
		reply, err := conn.Do("BLPOP", "blackidlist", int(time.Second))

		id, err := redis.Int(reply, err)

		if err != nil {
			conn.Close()
			continue
		}

		secKillConf.RWBlackLock.Lock()

		secKillConf.IdBlackMap[id] = true

		secKillConf.RWBlackLock.Unlock()
		conn.Close()

		logs.Info("sync id list from redis succ id[%v]", id)
	}

}

func initProxy2LayerRedis() (err error) {

	secKillConf.proxy2LayerRedisPool = &redis.Pool{MaxActive: secKillConf.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout: time.Duration(secKillConf.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
		MaxIdle: secKillConf.RedisProxy2LayerConf.RedisMaxIdle,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", secKillConf.RedisProxy2LayerConf.RedisAddr)
		}}

	conn := secKillConf.proxy2LayerRedisPool.Get()

	defer conn.Close()

	_, err = conn.Do("ping")

	if err != nil {
		logs.Error("ping proxy2layer redis failed")
		return
	}
	return
}
func initBlackRedis() (err error) {

	secKillConf.blackRedisPool = &redis.Pool{MaxIdle: secKillConf.RedisBlackConf.RedisMaxIdle, MaxActive: secKillConf.RedisBlackConf.RedisMaxActive, IdleTimeout: time.Duration(secKillConf.RedisBlackConf.RedisIdleTimeout) * time.Second, Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", secKillConf.RedisBlackConf.RedisAddr)
	}}

	conn := secKillConf.blackRedisPool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed err:%v", err)
		return
	}
	fmt.Println("init black redis success!!!")

	return
}

//加载黑名单
func loadBlackList() (err error) {

	secKillConf.IpBlackMap = make(map[string]bool, 10000)
	secKillConf.IdBlackMap = make(map[int]bool, 10000)

	err = initBlackRedis()
	if err != nil {
		err = fmt.Errorf("init black redis failed err:%v", err)
		return
	}
	conn := secKillConf.blackRedisPool.Get()
	defer conn.Close()
	//从hash 中读取全部的域或值
	reply, err := conn.Do("hgetall", "idblacklist")
	idList, err := redis.Strings(reply, err)

	if err != nil {
		logs.Warn("hget all failed err;%v", err)
		return
	}

	for _, v := range idList {
		secKillConf.IpBlackMap[v] = true
	}

	go SyncIpBlackList()
	go SyncIdBlackList()

	return
}

func IntentService(serviceConf *SecSkillConf) (err error) {
	secKillConf = serviceConf

	err = loadBlackList()
	if err != nil {
		err = fmt.Errorf("loadblacklist failed err:%v", err)
		return
	}

	logs.Debug("init IntentService succ")

	err = initProxy2LayerRedis()

	if err != nil {
		logs.Error("init proxy2layer redis failed err:%v", err)
		return
	}

	secKillConf.SecLimitMgr = &SecLimitMgr{UserLimitMap: make(map[int]*Limit), IpLimitMap: make(map[string]*Limit),}
	//被秒杀的商品
	secKillConf.SecReqChan = make(chan *SecRequest, 1000)

	//秒杀被确认的结果
	secKillConf.UserConnMap = make(map[string]chan *SecResult)

	initRedisProcessFunc()

	return
}
func initRedisProcessFunc() {
	for i := 0; i < secKillConf.WriteProxy2LayerGoroutineNum; i++ {
		go WriteHandle()
	}
	for i := 0; i < secKillConf.ReadProxy2LayerGoroutineNum; i++ {
		go ReadHandle()
	}
}
