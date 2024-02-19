package repository

import (
	"github.com/google/uuid"
	"my-app/module/user/domain"
)

type UserDTO struct {
	Id        uuid.UUID `gorm:"column:id;"`
	FirstName string    `gorm:"column:first_name;"`
	LastName  string    `gorm:"column:last_name;"`
	Email     string    `gorm:"column:email;"`
	Password  string    `gorm:"column:password;"`
	Salt      string    `gorm:"column:salt;"`
	Role      string    `gorm:"column:role;"`
}

func (dto *UserDTO) ToEntity() (*domain.User, error) {
	return domain.NewUser(dto.Id, dto.FirstName, dto.LastName, dto.Email, dto.Password, dto.Salt, domain.GetRole(dto.Role))
}
