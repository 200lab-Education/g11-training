package image

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"my-app/common"
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

func (r repo) Find(ctx context.Context, id uuid.UUID) (*common.Image, error) {
	var img Image

	if err := r.db.Table(TbName).Where("id = ?", id).First(&img).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return &common.Image{
		Id:              img.Id,
		Title:           img.Title,
		FileName:        img.FileName,
		FileUrl:         img.FileUrl,
		FileSize:        img.FileSize,
		FileType:        img.FileType,
		StorageProvider: img.StorageProvider,
		Status:          img.Status,
	}, nil
}

func (r repo) SetImageStatusActivated(ctx context.Context, id uuid.UUID) error {
	if err := r.db.Table(TbName).Where("id = ?", id).
		Updates(Image{Status: StatusActivated}).Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}
