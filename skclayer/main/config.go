package main

import (
	"goDep/skclayer/service"
	"github.com/astaxie/beego/config"
	"fmt"
	"strings"
	"github.com/astaxie/beego/logs"
)

var (
	Conf *service.SecLayerConf
)

func initConfig(adapterName, fileName string) (err error) {
	cfg, err := config.NewConfig(adapterName, fileName)
	if err != nil {
		return
	}

	logPath := cfg.String("logs::log_path")
	if len(logPath) == 0 {
		err = fmt.Errorf("get log file path faile")
		return
	}
	logLevel := cfg.String("logs::log_level")

	if len(logLevel) == 0 {
		err = fmt.Errorf("get log level failed")
		return
	}

	Conf = new(service.SecLayerConf)
	Conf.LogLevel = logLevel
	Conf.LogPath = logPath

	etcdServiceAddr := cfg.String("etcd::etcd_service_addr")

	if len(etcdServiceAddr) == 0 {
		err = fmt.Errorf("get ectd service addr failed")
		return
	}

	etcdSecKeyPrefix := cfg.String("etcd::etcd_sec_key_prefix")

	if len(etcdSecKeyPrefix) == 0 {
		err = fmt.Errorf("get etcd sec key failed")
		return
	}

	etcdProductKey := cfg.String("etcd::etcd_product_key")
	if len(etcdProductKey) == 0 {
		err = fmt.Errorf("get etcd product key failed")
		return
	}

	etcdBlackListKey := cfg.String("etcd::etcd_black_list_key")
	if len(etcdBlackListKey) == 0 {
		err = fmt.Errorf("get etcd black list key failed")
		return
	}

	etcdTimeout, err := cfg.Int("etcd::etcd_timeout")
	if err != nil {
		err = fmt.Errorf("get etcd timeout failed,err:%v", err)
		return
	}
	if !strings.HasSuffix(etcdSecKeyPrefix, "/") {
		etcdSecKeyPrefix = etcdSecKeyPrefix + "/"
	}

	Conf.EtcdConfig.EtcdSecKeyPrefix = etcdSecKeyPrefix
	Conf.EtcdConfig.EtcdProductKey = fmt.Sprint(Conf.EtcdConfig.EtcdSecKeyPrefix, etcdProductKey)
	Conf.EtcdConfig.EtcdBlackListKey = fmt.Sprint(Conf.EtcdConfig.EtcdSecKeyPrefix, etcdBlackListKey)
	Conf.EtcdConfig.EtcdServiceAddr = etcdServiceAddr
	Conf.EtcdConfig.EtcdTimeOut = etcdTimeout

	//读取redis相关的配置
	Conf.Proxy2LayerRedis.RedisAddr = cfg.String("redis::redis_proxy2layer_addr")
	if len(Conf.Proxy2LayerRedis.RedisAddr) == 0 {
		logs.Error("read redis::redis_proxy2layer_addr failed")
		err = fmt.Errorf("read redis::redis_proxy2layer_addr failed")
		return
	}

	Conf.Proxy2LayerRedis.RedisQueueName = cfg.String("redis::redis_proxy2layer_queue_name")
	if len(Conf.Proxy2LayerRedis.RedisQueueName) == 0 {
		logs.Error("read redis::resis_proxy2layer_queue_name failed")
		err = fmt.Errorf("read redis::resis_proxy2layer_queue_name failed")
		return
	}

	Conf.Proxy2LayerRedis.RedisMaxIdle, err = cfg.Int("redis::redis_proxy2layer_idle")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_idle failed, err:%v", err)
		return
	}

	Conf.Proxy2LayerRedis.RedisIdleTimeout, err = cfg.Int("redis::redis_proxy2layer_idle_time")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_idle_timeout failed, err:%v", err)
		return
	}

	Conf.Proxy2LayerRedis.RedisMaxActive, err = cfg.Int("redis::redis_proxy2layer_active")
	if err != nil {
		logs.Error("read redis::redis_proxy2layer_active failed, err:%v", err)
		return
	}

	//读取redis layer2proxy相关的配置
	Conf.Layer2ProxyRedis.RedisAddr = cfg.String("redis::redis_layer2proxy_addr")
	if len(Conf.Proxy2LayerRedis.RedisAddr) == 0 {
		logs.Error("read redis::redis_layer2proxy_addr failed")
		err = fmt.Errorf("read redis::redis_layer2proxy_addr failed")
		return
	}

	Conf.Layer2ProxyRedis.RedisQueueName = cfg.String("redis::redis_layer2proxy_queue_name")
	if len(Conf.Proxy2LayerRedis.RedisQueueName) == 0 {
		logs.Error("read redis::redis_layer2proxy_queue_name failed")
		err = fmt.Errorf("read redis::redis_layer2proxy_queue_name failed")
		return
	}

	Conf.Layer2ProxyRedis.RedisMaxIdle, err = cfg.Int("redis::redis_layer2proxy_idle")
	if err != nil {
		logs.Error("read redis::redis_layer2proxy_idle failed, err:%v", err)
		return
	}

	Conf.Layer2ProxyRedis.RedisIdleTimeout, err = cfg.Int("redis::redis_layer2proxy_idle_time")
	if err != nil {
		logs.Error("read redis::redis_layer2proxy_idle_timeout failed, err:%v", err)
		return
	}

	Conf.Layer2ProxyRedis.RedisMaxActive, err = cfg.Int("redis::redis_layer2proxy_active")
	if err != nil {
		logs.Error("read redis::redis_layer2proxy_active failed, err:%v", err)
		return
	}

	//读取各类goroutine线程数量
	Conf.ReadGoroutineNum, err = cfg.Int("service::read_layer2proxy_goroutine_num")
	if err != nil {
		logs.Error("read service::read_layer2proxy_goroutine_num failed, err:%v", err)
		return
	}

	Conf.WriteGoroutineNum, err = cfg.Int("service::write_proxy2layer_goroutine_num")
	if err != nil {
		logs.Error("read service::write_proxy2layer_goroutine_num failed, err:%v", err)
		return
	}

	Conf.HandleUserGoroutineNum, err = cfg.Int("service::handle_user_goroutine_num")
	if err != nil {
		logs.Error("read service::handle_user_goroutine_num failed, err:%v", err)
		return
	}

	Conf.Read2handleChanSize, err = cfg.Int("service::read2handle_chan_size")
	if err != nil {
		logs.Error("read service::read2handle_chan_size failed, err:%v", err)
		return
	}

	Conf.MaxRequestWaitTimeout, err = cfg.Int("service::max_request_wait_timeout")
	if err != nil {
		logs.Error("read service::max_request_wait_timeout failed, err:%v", err)
		return
	}

	Conf.Handle2WriteChanSize, err = cfg.Int("service::handle2write_chan_size")
	if err != nil {
		logs.Error("read service::handle2write_chan_size failed, err:%v", err)
		return
	}

	Conf.SendToWriteChanTimeout, err = cfg.Int("service::send_to_write_chan_timeout")
	if err != nil {
		logs.Error("read service::send_to_write_chan_timeout failed, err:%v", err)
		return
	}

	Conf.SendToHandleChanTimeout, err = cfg.Int("service::send_to_handle_chan_timeout")
	if err != nil {
		logs.Error("read service::send_to_handle_chan_timeout failed, err:%v", err)
		return
	}

	//读取token秘钥
	Conf.TokenPasswd = cfg.String("service::seckill_token_passwd")
	if len(Conf.TokenPasswd) == 0 {
		logs.Error("read service::seckill_token_passwd failed")
		err = fmt.Errorf("read service::seckill_token_passwd failed")
		return
	}

	return
}
