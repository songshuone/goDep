package router

import (
	"github.com/astaxie/beego"
	"goDep/secadmin/controller/product"
	"goDep/secadmin/controller/activity"
)

func init() {
	beego.Router("/product/submit", &product.ProductController{}, "*:SubmitProduct")
	beego.Router("/product/list", &product.ProductController{}, "*:ListProduct")
	beego.Router("/app/list", &product.ProductController{}, "*:ListProduct")
	beego.Router("/", &product.ProductController{}, "*:ListProduct")
	beego.Router("/product/create", &product.ProductController{}, "*:CreateProduct")

	beego.Router("/activity/list",&activity.ActivityController{},"*:ListActivity")
	beego.Router("/activity/create",&activity.ActivityController{},"*:CreateActivity")
	beego.Router("/activity/submit",&activity.ActivityController{},"*:SubmitActivity")

}
