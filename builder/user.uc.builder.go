package builder

import (
	"context"
	"gorm.io/gorm"
	"my-app/common"
	"my-app/module/user/domain"
	"my-app/module/user/infras/repository"
	"my-app/module/user/usecase"
)

type simpleBuilder struct {
	db *gorm.DB
	tp usecase.TokenProvider
}

func NewSimpleBuilder(db *gorm.DB, tp usecase.TokenProvider) simpleBuilder {
	return simpleBuilder{db: db, tp: tp}
}

func (s simpleBuilder) BuildUserQueryRepo() usecase.UserQueryRepository {
	return repository.NewUserRepo(s.db)
}

func (s simpleBuilder) BuildUserCmdRepo() usecase.UserCommandRepository {
	return repository.NewUserRepo(s.db)
}

func (simpleBuilder) BuildHasher() usecase.Hasher {
	return &common.Hasher{}
}

func (s simpleBuilder) BuildTokenProvider() usecase.TokenProvider {
	return s.tp
}

func (s simpleBuilder) BuildSessionQueryRepo() usecase.SessionQueryRepository {
	return repository.NewSessionMySQLRepo(s.db)
}

func (s simpleBuilder) BuildSessionCmdRepo() usecase.SessionCommandRepository {
	return repository.NewSessionMySQLRepo(s.db)
}

// Complex builder

func NewComplexBuilder(simpleBuilder simpleBuilder) complexBuilder {
	return complexBuilder{simpleBuilder: simpleBuilder}
}

type complexBuilder struct {
	simpleBuilder
}

// Proxy design pattern
type userCacheRepo struct {
	realRepo usecase.UserQueryRepository
	cache    map[string]*domain.User
}

func (c userCacheRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if user, ok := c.cache[email]; ok {
		return user, nil
	}

	user, err := c.realRepo.FindByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	c.cache[email] = user

	return user, nil
}

func (cb complexBuilder) BuildUserQueryRepo() usecase.UserQueryRepository {
	return userCacheRepo{
		realRepo: cb.simpleBuilder.BuildUserQueryRepo(),
		cache:    make(map[string]*domain.User),
	}
}
