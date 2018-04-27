package activity

import (
	"github.com/astaxie/beego"
	"goDep/secadmin/model"
	"github.com/astaxie/beego/logs"
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

	activity_name := a.GetString("activity_name")
	product_id, err := a.GetInt("product_id")
	if err != nil {
		logs.Warn("productid", product_id)
		return
	}
	start_time, err := a.GetInt64("start_time")
	if err != nil {
		logs.Warn("start_time", start_time)
		return
	}
	end_time, err := a.GetInt64("end_time")
	if err != nil {
		logs.Warn("end_time", end_time)
		return
	}
	total, err := a.GetInt("total")
	if err != nil {
		logs.Warn("total", total)
		return
	}
	speed, err := a.GetFloat("speed")
	if err != nil {
		logs.Warn("speed", speed)
		return
	}
	buy_limit, err := a.GetInt("buy_limit")
	if err != nil {
		logs.Warn("buy_limit", buy_limit)
		return
	}

	activityModel := model.NewActivityModel()
	activity := model.Activity{
	}

	activity.ProductId=product_id
	activity.StartTime=start_time
	activity.Total=total
	activity.EndTime=end_time
	activity.ActivityName=activity_name
	activity.BuyLimit=buy_limit
	activity.Speed=speed

	activityModel.CreateActivity(&activity)
	a.TplName = "activity/create.html"
	a.Layout = "layout/layout.html"
}
