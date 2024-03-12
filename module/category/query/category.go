package query

import (
	"context"
	"github.com/google/uuid"
	sctx "github.com/viettranx/service-context"
	"github.com/viettranx/service-context/core"
	"my-app/common"
)

type CategoryDTO struct {
	Id    uuid.UUID `gorm:"column:id" json:"id"`
	Title string    `gorm:"column:title" json:"title"`
}

func (CategoryDTO) TableName() string { return "categories" }

type categoriesByIdsQuery struct {
	sctx sctx.ServiceContext
}

func NewCategoriesByIdsQuery(sctx sctx.ServiceContext) *categoriesByIdsQuery {
	return &categoriesByIdsQuery{sctx: sctx}
}

func (q *categoriesByIdsQuery) Execute(ctx context.Context, ids []uuid.UUID) ([]CategoryDTO, error) {
	var cats []CategoryDTO

	dbContext := q.sctx.MustGet(common.KeyGorm).(common.DbContext)

	if err := dbContext.GetDB().Table(CategoryDTO{}.TableName()).
		Where("id in (?)", ids).
		Find(&cats).Error; err != nil {
		return nil, core.ErrBadRequest.WithError("cannot list categories").WithDebug(err.Error())
	}

	return cats, nil
}
