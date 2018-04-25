package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	_ "goDep/secadmin/router"
)

func main() {
	err := initAll()
	if err != nil {
		logs.Error(err)
		panic(err)
		return
	}
	beego.Run()
}
