package query

import (
	"context"
	"github.com/google/uuid"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"my-app/common"
)

// View/Read Model

type ProductDTO struct {
	common.BaseModel
	CatId    uuid.UUID    `gorm:"column:cat_id" json:"category_id"`
	Name     string       `gorm:"column:name" json:"name"`
	Type     string       `gorm:"column:type" json:"type"`
	Category *CategoryDTO `gorm:"foreignKey:cat_id;" json:"category"`
	//Description string      `gorm:"column:description" json:"description"`
}

type CategoryDTO struct {
	Id    uuid.UUID `gorm:"column:id" json:"id"`
	Title string    `gorm:"column:title" json:"title"`
}

func (CategoryDTO) TableName() string { return "categories" }

type ListProductFilter struct {
	CategoryId string `form:"category_id" json:"category_id"`
}

type ListProductParam struct {
	common.Paging
	ListProductFilter
}

type listProductQuery struct {
	sctx sctx.ServiceContext
}

func NewListProductQuery(sctx sctx.ServiceContext) listProductQuery {
	return listProductQuery{sctx: sctx}
}

func (q listProductQuery) Execute(ctx context.Context, param *ListProductParam) ([]ProductDTO, error) {
	var products []ProductDTO

	dbContext := q.sctx.MustGet(common.KeyGorm).(common.DbContext)

	db := dbContext.GetDB().Table("products")

	if param.CategoryId != "" {
		db = db.Where("category_id = ?", param.CategoryId)
	}

	var count int64
	db.Count(&count)
	param.Total = int(count)

	param.Process()

	offset := param.Limit * (param.Page - 1)

	db = db.Preload("Category")

	if err := db.Offset(offset).Limit(param.Limit).Order("id desc").Find(&products).Error; err != nil {
		return nil, core.ErrBadRequest.WithError("cannot list product").WithDebug(err.Error())
	}

	return products, nil
}

type CategoryRepository interface {
	FindWithIds(ctx context.Context, ids []uuid.UUID) ([]CategoryDTO, error)
}
