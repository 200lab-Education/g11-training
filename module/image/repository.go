package image

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) repo {
	return repo{db: db}
}

func (r repo) Create(ctx context.Context, entity *Image) error {
	if err := r.db.Table(TbName).Create(entity).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
