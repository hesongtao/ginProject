package product

import (
	"fmt"
	"ginProject/common/dao"
	"gorm.io/gorm"
)

type Product struct {
	ID   int
	Name string
}

// 因为要和数据库建立连接，所以我们要把链接保存到var中
var (
	Db  *gorm.DB
	err error
)

func CreateProduct(name string) {
	// Read
	dao.Db.Create(&Product{ID: 30, Name: name})

}

func GetProductsById(id int) string {
	var product Product
	dao.Db.First(&product, id)
	if dao.Db.Error != nil {
		fmt.Println("-----database err:", err)
	}
	return product.Name

}
