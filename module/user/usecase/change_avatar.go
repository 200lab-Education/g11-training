package usecase

import (
	"context"
)

type changeAvtUC struct {
	userQueryRepo UserQueryRepository
	userCmdRepo   UserCommandRepository
	imgRepo       ImageRepository
}

func NewChangeAvtUC(userQueryRepo UserQueryRepository, userCmdRepo UserCommandRepository, imgRepo ImageRepository) *changeAvtUC {
	return &changeAvtUC{userQueryRepo: userQueryRepo, userCmdRepo: userCmdRepo, imgRepo: imgRepo}
}

func (uc *changeAvtUC) ChangeAvatar(ctx context.Context, dto SingleImageDTO) error {
	return nil
}
