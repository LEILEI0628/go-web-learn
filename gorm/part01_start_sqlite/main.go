package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("gorm/test.db"), &gorm.Config{
		//DryRun: true, // 只输出语句不会执行
	})
	if err != nil {
		panic("failed to connect database")
	}

	db = db.Debug()
	// db.Debug().Method debug模式运行，打印出SQL语句（更推荐通过logger实现）

	// 迁移 schema（建表）
	db.AutoMigrate(&Product{})

	// Create（插入）
	db.Create(&Product{Code: "D42", Price: 100})

	// Read（搜索）
	var product Product
	db.First(&product, 1)                 // 根据整型主键查找（传一个参数时，主键==1）
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录（传两个参数时）

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200) // 注意GORM框架不会对值的类型判断或转换
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&product, 1)
}
