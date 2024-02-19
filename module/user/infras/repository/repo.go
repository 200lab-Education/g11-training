package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"my-app/common"
	"my-app/module/user/domain"
)

const TbName = "users"

type userMySQLRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) userMySQLRepo {
	return userMySQLRepo{db: db}
}

func (repo userMySQLRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dto UserDTO

	if err := repo.db.Table(TbName).Where("email = ?", email).First(&dto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}

		return nil, err
	}

	return dto.ToEntity()
}

func (repo userMySQLRepo) Create(ctx context.Context, data *domain.User) error {
	dto := UserDTO{
		Id:        data.Id(),
		FirstName: data.FirstName(),
		LastName:  data.LastName(),
		Email:     data.Email(),
		Password:  data.Password(),
		Salt:      data.Salt(),
		Role:      data.Role().String(),
	}

	if err := repo.db.Table(TbName).Create(&dto).Error; err != nil {
		return err
	}

	return nil
}
