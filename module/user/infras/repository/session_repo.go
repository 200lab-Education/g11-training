package repository

import (
	"context"
	"gorm.io/gorm"
	"my-app/module/user/domain"
)

const TbSessionName = "user_sessions"

type sessionMySQLRepo struct {
	db *gorm.DB
}

func NewSessionMySQLRepo(db *gorm.DB) sessionMySQLRepo {
	return sessionMySQLRepo{db: db}
}

func (repo sessionMySQLRepo) Create(ctx context.Context, data *domain.Session) error {
	dto := SessionDTO{
		Id:           data.Id(),
		UserId:       data.UserId(),
		RefreshToken: data.RefreshToken(),
		AccessExpAt:  data.AccessExpAt(),
		RefreshExpAt: data.RefreshExpAt(),
	}

	return repo.db.Table(TbSessionName).Create(&dto).Error
}
