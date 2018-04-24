package service

import (
	"sync"

	"fmt"
	"github.com/astaxie/beego/logs"
)

//限制恶意用户

type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	lock         sync.Mutex
}

func antiSpam(req *SecRequest) (err error) {
	_, ok := secKillConf.IdBlackMap[req.UserId]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Warn("userid[%d] is block by id black", req.UserId)
		return
	}
	_, ok = secKillConf.IpBlackMap[req.ClientAddr]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Warn("client ip  [%d] is block by id black", req.ClientAddr)
		return
	}
	//ip频率控制
	ipLimit, ok := secKillConf.SecLimitMgr.IpLimitMap[req.ClientAddr]
	if !ok {
		ipLimit = &Limit{secLimit: &SecLimit{}, minLimit: &MinLimit{},}
		secKillConf.SecLimitMgr.IpLimitMap[req.ClientAddr] = ipLimit
	}
	secKillConf.SecLimitMgr.lock.Lock()
	ipMinLimitAccessCount := ipLimit.minLimit.Count(req.AccessTime.Unix())
	ipSecLimitAccessCount := ipLimit.secLimit.Count(req.AccessTime.Unix())

	//uid 频率控制
	idLimit, ok := secKillConf.SecLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		idLimit = &Limit{secLimit: &SecLimit{}, minLimit: &MinLimit{}}
		secKillConf.SecLimitMgr.UserLimitMap[req.UserId] = idLimit
	}
	idMinLimitAccessCount := ipLimit.minLimit.Count(req.AccessTime.Unix())
	idSecLimitAccessCount := ipLimit.secLimit.Count(req.AccessTime.Unix())
	secKillConf.SecLimitMgr.lock.Unlock()

	if ipMinLimitAccessCount > secKillConf.IpMinAccessLimit {
		err = fmt.Errorf("invalid request")
		secKillConf.IpBlackMap[req.ClientAddr] = true
		return
	}
	if ipSecLimitAccessCount > secKillConf.IpSecAccessLimit {
		err = fmt.Errorf("invalid request")
		secKillConf.IpBlackMap[req.ClientAddr] = true
		return
	}

	if idMinLimitAccessCount > secKillConf.IpMinAccessLimit {
		err = fmt.Errorf("invalid request")
		secKillConf.IdBlackMap[req.UserId] = true
		return
	}
	if idSecLimitAccessCount > secKillConf.IpSecAccessLimit {
		err = fmt.Errorf("invalid request")
		secKillConf.IdBlackMap[req.UserId] = true
		return
	}

	return
}
