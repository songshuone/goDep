package controller

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"goDep/skcproxy/service"
	"strings"
	"strconv"
	"time"
)

type SkillController struct {
	beego.Controller
}

//秒杀
func (s *SkillController) SecKill() {

	productId, err := s.GetInt("product_id")

	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "success"
	defer func() {
		s.Data["json"] = result
		s.ServeJSON()
	}()
	if err != nil {
		result["code"] = 1001
		result["message"] = "invalid product_id"
		return
	}

	source := s.GetString("src") //来源
	authCode := s.GetString("authcode")
	secTime := s.GetString("time") //时间
	nance := s.GetString("nance")  //随机数

	secRequest := service.NewSecRequest()

	secRequest.Source = source
	secRequest.AuthCode = authCode
	secRequest.SecTime = secTime
	secRequest.Nance = nance
	secRequest.ProductId = productId

	secRequest.UserAuthSign = s.Ctx.GetCookie("userAuthSign")
	secRequest.UserId, err = strconv.Atoi(s.Ctx.GetCookie("userId"))

	secRequest.AccessTime = time.Now()

	if len(s.Ctx.Request.RemoteAddr) > 0 {
		//获取客户端的ip地址
		secRequest.ClientAddr = strings.Split(s.Ctx.Request.RemoteAddr, ":")[0]
	}
	//获取客户端从那个入口访问的接口
	secRequest.ClientRefence = s.Ctx.Request.Referer()
	//获取是否客户端关闭连接
	secRequest.CloseNotify = s.Ctx.ResponseWriter.CloseNotify()

	logs.Debug("client request:[%v]", secRequest)

	if err != nil {
		result["code"] = service.ErrInvalidRequest
		result["message"] = "invalid cookie:userId"
		return
	}

	data, code, err := service.SecKill(secRequest)

	if err != nil {
		result["code"] = code
		result["message"] = err.Error()
		return
	}
	result["data"] = data
	result["code"] = code
	return
}

//秒杀信息
func (s *SkillController) SecInfo() {

	productId, err := s.GetInt("product_id")
	result := make(map[string]interface{})
	defer func() {
		s.Data["json"] = result
		s.ServeJSON()
	}()

	result["code"] = 0
	result["message"] = "success"

	if err != nil {
		data, _, err := service.SecInfo(productId)
		if err != nil {
			result["code"] = 1001
			result["message"] = "invalid product_id"
			logs.Error("invalid request,get product_id failed,err:%v", err)
			return
		}
		result["data"] = data
		return
	}
	data, code, err := service.SecInfo(productId)
	if err != nil {
		result["code"] = code
		result["message"] = err.Error()
		logs.Error("invalid request,get product_id failed,err:%v", err)
		return
	}
	result["data"] = data


}
