package service

import (
	"sync"
	"github.com/garyburd/redigo/redis"
	etcd_client "github.com/coreos/etcd/clientv3"
	"time"
)

var (
	secLayerContext = &SecLayerContext{}
)

//秒杀配置
type SecProductInfoConf struct {
	ProductId         int `json:"product_id"`
	StartTime         int64 `json:"start_time"`
	EndTime           int64 `json:"end_time"`
	Status            int `json:"status"`
	Total             int `json:"total"`
	Left              int `json:"left"`
	OnePersonBuyLimit int `json:"one_person_buy_limit"`
	BuyRate           float64 `json:"buy_rate"`
	//每秒最多能卖多少个
	SoldMaxLimit int `json:"sold_max_limit"`
	//限速控制
	secLimit *SecLimit
}

type SecLayerConf struct {
	LogLevel         string
	LogPath          string
	EtcdConfig       EtcdConfig
	Proxy2LayerRedis RedisConf
	Layer2ProxyRedis RedisConf

	WriteGoroutineNum      int
	ReadGoroutineNum       int
	HandleUserGoroutineNum int
	Read2handleChanSize    int
	Handle2WriteChanSize   int
	MaxRequestWaitTimeout  int

	SendToWriteChanTimeout  int
	SendToHandleChanTimeout int
	TokenPasswd             string

	SecProductInfoMap map[int]*SecProductInfoConf
}

//etcd_service_addr=127.0.0.1:2379
//etcd_sec_key_prefix=/backend/secskill
//etcd_product_key=product
//etcd_black_list_key=backlist
//etcd_timeout=10
type EtcdConfig struct {
	EtcdServiceAddr  string
	EtcdSecKeyPrefix string
	EtcdProductKey   string
	EtcdBlackListKey string
	EtcdTimeOut      int
}

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
	RedisQueueName   string
}

type SecLayerContext struct {
	proxy2LayerRedisPool *redis.Pool
	layer2ProxyRedisPool *redis.Pool
	etcdClient           *etcd_client.Client
	RWSecProductLock     sync.RWMutex

	secLayerConf *SecLayerConf
	waitGroup    sync.WaitGroup
	Read2HandleChan  chan *SecRequest
	Handle2WriteChan chan *SecResponse

	HistoryMap     map[int]*UserBuyHistory
	HistoryMapLock sync.Mutex

	//商品的计数
	productCountMgr *ProductCountMgr
}
type SecRequest struct {
	ProductId     int
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time
	ClientAddr    string
	ClientRefence string
	//closeNotify   <-chan bool

	//ResultChan chan *SecResult
}

type SecResponse struct {
	ProductId int `json:"product_id"`
	UserId    int `json:"user_id"`
	Token     string `json:"token"`
	TokenTime int64 `json:"token_time"`
	Code      int `json:"code"`
}