package product

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"fmt"
	"goDep/secadmin/model"
)

type ProductController struct {
	beego.Controller
}

func (p *ProductController) CreateProduct() {
	p.TplName = "product/create.html"
	p.Layout = "layout/layout.html"
}

func (p *ProductController) ListProduct() {
	productModel := model.NewProductModel()

	productList, err := productModel.GetProductList()
	if err != nil {
		logs.Warn("get product list failed,err:%v", err)
		return
	}
	p.Data["product_list"] = productList
	p.TplName = "product/list.html"
	p.Layout = "layout/layout.html"
}

func (c *ProductController) SubmitProduct() {

	productName := c.GetString("product_name")
	productTotal, err := c.GetInt("product_total")
	c.TplName = "product/create.html"
	c.Layout = "layout/layout.html"
	errMsg := "success"

	defer func() {
		if err != nil {
			c.Data["Error"] = errMsg
			c.TplName = "product/error.html"
			c.Layout = "layout/layout.html"
		}
	}()

	if len(productName) == 0 {
		logs.Warn("invalid product name,err:%v", err)
		errMsg = fmt.Sprintf("invalid product name,err:%v", err)
		return
	}
	if err != nil {
		logs.Warn("invalid product total,err:%v", err)
		errMsg = fmt.Sprintf("invalid product total,err:%v", err)
		return
	}

	productStatus, err := c.GetInt("product_status")

	if err != nil {
		logs.Warn("invalid product status,err:%v", err)
		return
	}

	productModel := model.NewProductModel()

	product := model.Product{
		ProductName: productName,
		Total:       productTotal,
		Status:      productStatus,
	}

	err = productModel.CreateProduct(&product)

	if err != nil {
		logs.Warn("create product failed,err:%v", err)
		errMsg = fmt.Sprintf("create product failed,err:%v", err)
		return
	}
	logs.Debug("product name[%s] product total[%d],product status[%d]", productName, productTotal, productStatus)

}
