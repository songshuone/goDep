package activity

import (
	"github.com/astaxie/beego"
	"goDep/secadmin/model"
	"github.com/astaxie/beego/logs"
	"net/http"
)

type ActivityController struct {
	beego.Controller
}

func (a *ActivityController) ListActivity() {
	activityModel := model.NewActivityModel()
	activityList, err := activityModel.GetActivityList()
	if err != nil {
		return
	}
	a.Data["activity_list"] = activityList
	a.TplName = "activity/list.html"
	a.Layout = "layout/layout.html"
}

func (a *ActivityController) CreateActivity() {
	a.TplName = "activity/create.html"
	a.Layout = "layout/layout.html"
}

func (a *ActivityController) SubmitActivity() {

	var Error string = "success"
	a.Layout = "layout/layout.html"
	defer func() {
		a.Data["Error"] = Error
		a.TplName = "activity/error.html"
	}()
	activity_name := a.GetString("activity_name")
	if len(activity_name) == 0 {
		logs.Warn("activity_name", activity_name)
		Error = "活动的名字不能为空字符串"
		return
	}
	product_id, err := a.GetInt("product_id")

	if err != nil {
		logs.Warn("productid", product_id)
		Error = "活动的商品Id非法"
		return
	}
	start_time, err := a.GetInt64("start_time")
	if err != nil {
		logs.Warn("start_time", start_time)
		Error = "活动的开始时间非法"
		return
	}
	end_time, err := a.GetInt64("end_time")
	if err != nil {
		Error = "活动的结束时间非法"
		logs.Warn("end_time", end_time)
		return
	}
	total, err := a.GetInt("total")
	if err != nil {
		Error = "活动的总数非法"
		logs.Warn("total", total)
		return
	}
	speed, err := a.GetFloat("speed")
	if err != nil {
		Error = "活动的商品速度非法"
		logs.Warn("speed", speed)
		return
	}
	buy_limit, err := a.GetInt("buy_limit")
	if err != nil {
		logs.Warn("buy_limit", buy_limit)
		Error = "活动的限制数量非法"
		return
	}

	activityModel := model.NewActivityModel()
	activity := model.Activity{
	}

	activity.ProductId = product_id
	activity.StartTime = start_time
	activity.Total = total
	activity.EndTime = end_time
	activity.ActivityName = activity_name
	activity.BuyLimit = buy_limit
	activity.Speed = speed

	err = activityModel.CreateActivity(&activity)
	if err != nil {
		logs.Warn("create Activity failed,err:%v", err)
		Error = "create Activity failed,err:" + err.Error()
		return
	}
	a.Redirect("/activity/list", http.StatusMovedPermanently)

}
