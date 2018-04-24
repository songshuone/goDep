package main

import (
	"github.com/astaxie/beego"
	"fmt"
	"strings"
	"goDep/skcproxy/service"
)

func initConfig() (err error) {

	redisBlackAddr := beego.AppConfig.String("redis_black_addr")
	redisBlackIdle, err := beego.AppConfig.Int("redis_black_idle")

	redisBlackctive, err := beego.AppConfig.Int("redis_black_active")
	redisBlackIdleTimeout, err := beego.AppConfig.Int("redis_black_idle_timeout")

	if len(redisBlackAddr) == 0 {
		err = fmt.Errorf("redis addr is  null")
		return
	}
	if err != nil {
		err = fmt.Errorf("redis config failed err:%v", err)
		return
	}
	service.Conf.RedisBlackConf.RedisAddr = redisBlackAddr
	service.Conf.RedisBlackConf.RedisMaxIdle = redisBlackIdle
	service.Conf.RedisBlackConf.RedisIdleTimeout = redisBlackIdleTimeout
	service.Conf.RedisBlackConf.RedisMaxActive = redisBlackctive

	redisProxy2LayerAddr := beego.AppConfig.String("redis_proxy2layer_addr")
	redisProxy2LayerIdle, err := beego.AppConfig.Int("redis_proxy2layer_idle")
	redisProxy2LayerActive, err := beego.AppConfig.Int("redis_proxy2layer_active")
	redisProxy2LayerIdleTimeout, err := beego.AppConfig.Int("redis_proxy2layer_idle_timeout")
	if len(redisProxy2LayerAddr) == 0 {
		err = fmt.Errorf("redisproxy2layeraddr:%v", redisProxy2LayerAddr)
		return
	}
	if err != nil {
		err = fmt.Errorf("init redisproxy2layer failed err:%v", err)
		return
	}
	service.Conf.RedisProxy2LayerConf.RedisAddr = redisProxy2LayerAddr
	service.Conf.RedisProxy2LayerConf.RedisMaxActive = redisProxy2LayerActive
	service.Conf.RedisProxy2LayerConf.RedisIdleTimeout = redisProxy2LayerIdleTimeout
	service.Conf.RedisProxy2LayerConf.RedisMaxIdle = redisProxy2LayerIdle

	redisLayer2ProxyAddr := beego.AppConfig.String("redis_layer2proxy_addr")
	redisLayer2ProxyIdle, err := beego.AppConfig.Int("redis_layer2proxy_idle")
	redisLayer2ProxyActive, err := beego.AppConfig.Int("redis_layer2proxy_active")
	redisLayer2ProxyIdleTimeout, err := beego.AppConfig.Int("redis_layer2proxy_idle_timeout")
	if len(redisProxy2LayerAddr) == 0 {
		err = fmt.Errorf("redislayer2proxyaddr:%v", redisProxy2LayerAddr)
		return
	}
	if err != nil {
		err = fmt.Errorf("init redislayer2proxy failed err:%v", err)
		return
	}
	service.Conf.RedisLayer2ProxyConf.RedisAddr = redisLayer2ProxyAddr
	service.Conf.RedisLayer2ProxyConf.RedisMaxActive = redisLayer2ProxyActive
	service.Conf.RedisLayer2ProxyConf.RedisIdleTimeout = redisLayer2ProxyIdleTimeout
	service.Conf.RedisLayer2ProxyConf.RedisMaxIdle = redisLayer2ProxyIdle

	etcdAddr := beego.AppConfig.String("etcd_addr")
	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")
	etcdSecKeyPrefix := beego.AppConfig.String("etcd_sec_key_prefix")
	etcdProductKey := beego.AppConfig.String("etcd_product_key")
	etcdBlankListKey := beego.AppConfig.String("etcd_black_list_key")

	if err != nil {
		return
	}
	if len(etcdAddr) == 0 || len(etcdSecKeyPrefix) == 0 || len(etcdProductKey) == 0 || len(etcdBlankListKey) == 0 {
		err = fmt.Errorf("etcd get config  fail")
		return
	}

	if !strings.HasSuffix(etcdSecKeyPrefix, "/") {
		etcdSecKeyPrefix = fmt.Sprint(etcdSecKeyPrefix, "/")
	}

	service.Conf.EtcdAddr = etcdAddr
	service.Conf.EtcdTimeout = etcdTimeout
	service.Conf.EtcdSecKeyPrefix = etcdSecKeyPrefix
	service.Conf.EtcdBlankListKey = etcdBlankListKey
	service.Conf.EtcdProductKey = fmt.Sprint(etcdSecKeyPrefix, etcdProductKey)

	logPath := beego.AppConfig.String("log_path")
	logLevel := beego.AppConfig.String("log_level")

	if len(logPath) == 0 || len(logLevel) == 0 {
		err = fmt.Errorf("logpath:[%v],loglevel:[%v]", logPath, logLevel)
		return
	}
	service.Conf.LogPath = logPath
	service.Conf.LogLevel = logLevel

	referWhitList := beego.AppConfig.String("refer_white_list")
	if len(referWhitList) > 0 {
		service.Conf.ReferWhiteList = strings.Split(referWhitList, ",")
	}

	cookieSecretKey := beego.AppConfig.String("cookie_secret_key")
	if len(cookieSecretKey) > 0 {
		service.Conf.CookieSecretKey = cookieSecretKey
	}

	//获取频率控制阈值
	ipSecAccessLimit, err := beego.AppConfig.Int("ip_sec_access_limit")
	ipMinAccessLimit, err := beego.AppConfig.Int("ip_min_access_limit")
	userSecAccessLimit, err := beego.AppConfig.Int("user_sec_access_limit")
	userMinAccessLimit, err := beego.AppConfig.Int("user_min_access_limit")

	if err != nil {
		err = fmt.Errorf("init config failed err：[%v] ,ip_sec_access_limit|ip_min_access_limit|user_sec_access_limit|user_min_access_limit", err)
		return
	}
	service.Conf.IpSecAccessLimit = ipSecAccessLimit
	service.Conf.IpMinAccessLimit = ipMinAccessLimit
	service.Conf.UserMinAccessLimit = userMinAccessLimit
	service.Conf.UserSecAccessLimit = userSecAccessLimit

	writeProxy2LayerGoroutinNunm, err := beego.AppConfig.Int("write_proxy2layer_goroutine_num")

	if err != nil {
		writeProxy2LayerGoroutinNunm = 10
	}

	readProxy2LayerGoroutineNum, err := beego.AppConfig.Int("read_proxy2layer_goroutine_num")
	if err != nil {
		readProxy2LayerGoroutineNum = 10
	}
	service.Conf.ReadProxy2LayerGoroutineNum = readProxy2LayerGoroutineNum
	service.Conf.WriteProxy2LayerGoroutineNum = writeProxy2LayerGoroutinNunm
	return
}
