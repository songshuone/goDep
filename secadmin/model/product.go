package model

import (
	"github.com/astaxie/beego/logs"
)

type ProductModel struct {
}
type Product struct {
	ID          int
	ProductName string
	Total       int
	Status      int
}

func NewProductModel() *ProductModel {
	return &ProductModel{}
}

func (p *ProductModel) CreateProduct(product *Product) (err error) {
	err = Db.Create(product).Error
	if err != nil {
		logs.Warn("create product failed,err:%v", err)
		return
	}
	logs.Debug("create product success")
	return
}

func (p *ProductModel) GetProductList() (product []*Product, err error) {
	err = Db.Find(&product).Error
	return
}
