package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"goDep/skclayer/service"
)

func main() {

	err := initConfig("ini", "./conf/seclayer.conf")
	if err != nil {
		panic(fmt.Sprintf("init config failed err:%v", err))
		return
	}

	err = initLogger()
	if err != nil {
		panic(fmt.Sprintf("init logger failed err:%v", err))
		return
	}
	logs.Info("init logger success!!!")
	b, _ := json.Marshal(Conf)
	logs.Info(string(b))

	err = service.InitSecLayer(Conf)
	if err != nil {
		logs.Error("init seclayer failed, err:%v", err)
		return
	}
	logs.Debug("init layer succss")

	err = service.Run()
	if err != nil {
		logs.Error("service run return, err:%v", err)
		return
	}
	logs.Info("service run exited")
}
