package service

import (
	"sync"
	"time"
	"github.com/garyburd/redigo/redis"
)

var (
	Conf = &SecSkillConf{
		SecProductInfoMap: make(map[int]*SecProductInfoConf, 1024),
	}
)

type SecSkillConf struct {
	RedisBlackConf       RedisConfig
	RedisProxy2LayerConf RedisConfig
	RedisLayer2ProxyConf RedisConfig
	EtcdConfig
	LogConfig
	SecProductInfoMap    map[int]*SecProductInfoConf
	RWSecProductLock     sync.RWMutex

	ReferWhiteList  []string
	CookieSecretKey string

	AccessLimitConf

	IpBlackMap map[string]bool
	IdBlackMap map[int]bool

	SecLimitMgr *SecLimitMgr

	UserConnMap     map[string]chan *SecResult
	SecReqChan      chan *SecRequest
	UserConnMapLock sync.Mutex

	blackRedisPool       *redis.Pool
	proxy2LayerRedisPool *redis.Pool
	layer2ProxyRedisPool *redis.Pool
	RWBlackLock sync.Mutex

	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int
}

type AccessLimitConf struct {
	IpSecAccessLimit   int
	IpMinAccessLimit   int
	UserSecAccessLimit int
	UserMinAccessLimit int
}

type RedisConfig struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type EtcdConfig struct {
	EtcdAddr         string
	EtcdTimeout      int
	EtcdSecKeyPrefix string
	EtcdProductKey   string
	EtcdBlankListKey string
}

type LogConfig struct {
	LogPath  string
	LogLevel string
}

//秒杀结果数据


type SecResult struct {
	ProductId int `json:"product_id"`
	UserId    int `json:"user_id"`
	Code      int `json:"code"`
	Token     string `json:"token"`
}

type SecRequest struct {
	ProductId     int
	Source        string
	AuthCode      string
	SecTime       string
	Nance         string
	UserId        int
	UserAuthSign  string
	AccessTime    time.Time //访问的时间
	ClientAddr    string    //记录客户端的地址
	ClientRefence string    //记录从什么地址来访问的

	CloseNotify <-chan bool `json:"_"`

	ResultChan chan *SecResult `json:"_"`
}
