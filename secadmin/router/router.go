package router

import (
	"github.com/astaxie/beego"
	"goDep/secadmin/controller/product"
)

func init() {
	beego.Router("/product/submit", &product.ProductController{}, "*:SubmitProduct")

}
