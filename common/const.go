package common

import "gorm.io/gorm"

const (
	KeyRequester = "requester"
	KeyGorm      = "gorm"
	KeyJWT       = "jwt"
)

type DbContext interface {
	GetDB() *gorm.DB
}
