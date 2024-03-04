package image

import (
	"context"
	"errors"
	"fmt"
	"github.com/viettranx/service-context/core"
	"my-app/common"
	"time"
)

type UseCase interface {
	UploadImage(ctx context.Context, dto UploadDTO) (*Image, error)
}

type useCase struct {
	uploader Uploader
	repo     CmdRepository
}

func NewUseCase(uploader Uploader, repo CmdRepository) useCase {
	return useCase{uploader: uploader, repo: repo}
}

func (uc useCase) UploadImage(ctx context.Context, dto UploadDTO) (*Image, error) {
	dstFileName := fmt.Sprintf("%d_%s", time.Now().UTC().UnixNano(), dto.FileName)

	if err := uc.uploader.SaveFileUploaded(ctx, dto.FileData, dstFileName); err != nil {
		return nil, core.ErrInternalServerError.WithError(ErrCannotUploadImage.Error()).WithDebug(err.Error())
	}

	image := Image{
		Id:              common.GenUUID(),
		Title:           dto.Name,
		FileName:        dstFileName,
		FileSize:        dto.FileSize,
		FileType:        dto.FileType,
		StorageProvider: uc.uploader.GetName(),
		Status:          StatusUploaded,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := uc.repo.Create(ctx, &image); err != nil {
		return nil, core.ErrInternalServerError.WithError(ErrCannotUploadImage.Error()).WithDebug(err.Error())
	}

	return &image, nil
}

type Uploader interface {
	SaveFileUploaded(ctx context.Context, data []byte, dst string) error
	GetName() string
	GetDomain() string
}

type CmdRepository interface {
	Create(ctx context.Context, entity *Image) error
}

type UploadDTO struct {
	Name     string
	FileName string
	FileType string
	FileSize int
	FileData []byte
}

var (
	ErrCannotUploadImage = errors.New("cannot upload image")
	ErrCannotFindImage   = errors.New("cannot find image")
)
