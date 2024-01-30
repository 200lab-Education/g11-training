package main

import (
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

//
//`id` int NOT NULL AUTO_INCREMENT,
//`category_id` int DEFAULT NULL,
//`name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
//`image` json DEFAULT NULL,
//`type` enum('drink','food','topping') NOT NULL DEFAULT 'drink',
//`description` text,
//`status` enum('activated','deactivated','out_of_stock') DEFAULT 'activated',
//`created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
//`updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

type BaseModel struct {
	Id        uuid.UUID `gorm:"column:id;"`
	Status    string    `gorm:"column:status;"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type Product struct {
	BaseModel
	CategoryId int    `gorm:"column:category_id"`
	Name       string `gorm:"column:name"`
	//Image       any    `gorm:"column:image"`
	Type        string `gorm:"column:type"`
	Description string `gorm:"column:description"`
}

type ProductUpdate struct {
	Name        *string `gorm:"column:name"`
	CategoryId  *int    `gorm:"column:category_id"`
	Status      *string `gorm:"column:status;"`
	Type        *string `gorm:"column:type"`
	Description *string `gorm:"column:description"`
}

func (Product) TableName() string { return "products" }

func main() {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	now := time.Now().UTC()

	newId, _ := uuid.NewV7()

	newProd := Product{
		BaseModel: BaseModel{
			Id:        newId,
			Status:    "activated",
			CreatedAt: now,
			UpdatedAt: now,
		},
		CategoryId: 1,
		Name:       "Latte",
		//Image:       nil,
		Type:        "drink",
		Description: "",
	}

	if err := db.Table("products").Create(&newProd).Error; err != nil {
		log.Println(err)
	}

	var oldProduct Product

	if err := db.
		Table(Product{}.TableName()).
		//Where("id = ?", 4).
		First(&oldProduct).Error; err != nil {
		log.Println(err)
	}

	log.Println("Product:", oldProduct)

	var prods []Product

	if err := db.
		Table(Product{}.TableName()).
		Where("status not in (?)", []string{"deactivated"}).
		Limit(10).
		Offset(10).
		Order("id desc").
		Find(&prods).Error; err != nil {
		log.Println(err)
	}

	log.Println("Products:", prods)

	//oldProduct.Name = ""

	//emptyStr := "Latte"

	//if err := db.
	//	Table(Product{}.TableName()).
	//	Where("id = ?", 4).
	//	Updates(ProductUpdate{Name: &emptyStr}).Error; err != nil {
	//	log.Println(err)
	//}

	//if err := db.
	//	Table(Product{}.TableName()).
	//	Where("id = ?", 4).
	//	Delete(nil).Error; err != nil {
	//	log.Println(err)
	//}
}
