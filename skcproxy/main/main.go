package main

import (
	"github.com/astaxie/beego"
	_ "goDep/skcproxy/router"
	"goDep/skcproxy/service"
)

func main() {

	err := initConfig()
	if err != nil {
		panic(err)
		return
	}

	beego.Info("init config success:", service.Conf)

	err = initSec()
	if err != nil {
		panic(err)
		return
	}
	beego.Run()
}
