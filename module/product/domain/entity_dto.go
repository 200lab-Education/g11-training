package productdomain

import "github.com/google/uuid"

type ProductUpdateDTO struct {
	Name        *string `gorm:"column:name" json:"name"`
	CategoryId  *int    `gorm:"column:category_id" json:"category_id"`
	Status      *string `gorm:"column:status;" json:"status"`
	Type        *string `gorm:"column:type" json:"type"`
	Description *string `gorm:"column:description" json:"description"`
}

type ProductCreationDTO struct {
	Id          uuid.UUID `gorm:"column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	CategoryId  int       `gorm:"column:category_id" json:"category_id"`
	Type        string    `gorm:"column:type" json:"type"`
	Description string    `gorm:"column:description" json:"description"`
}

func (ProductCreationDTO) TableName() string { return "products" }
func (ProductUpdateDTO) TableName() string   { return ProductCreationDTO{}.TableName() }
