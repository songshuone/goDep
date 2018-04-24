package router

import (
	"github.com/astaxie/beego"
	"goDep/skcproxy/controller"
)

func init() {
	beego.Router("/seckill", &controller.SkillController{}, "*:SecKill")
	beego.Router("/secinfo", &controller.SkillController{}, "*:SecInfo")
}
