package usecase

import (
	"context"
	"errors"
	"my-app/common"
	"my-app/module/user/domain"
)

type UseCase interface {
	Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error
}

type Hasher interface {
	RandomStr(length int) (string, error)
	HashPassword(salt, password string) (string, error)
}

type useCase struct {
	repo   UserRepository
	hasher Hasher
}

func NewUseCase(repo UserRepository, hasher Hasher) UseCase {
	return &useCase{repo: repo, hasher: hasher}
}

func (uc *useCase) Register(ctx context.Context, dto EmailPasswordRegistrationDTO) error {
	// 1. Find user by email:
	// 1.1 Found: return error (email has existed)
	// 2. Generate salt
	// 3. Hash password+salt
	// 4. Create user entity

	user, err := uc.repo.FindByEmail(ctx, dto.Email)

	if user != nil {
		return domain.ErrEmailHasExisted
	}

	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		return err
	}

	salt, err := uc.hasher.RandomStr(30)

	if err != nil {
		return err
	}

	hashedPassword, err := uc.hasher.HashPassword(salt, dto.Password)

	if err != nil {
		return err
	}

	userEntity, err := domain.NewUser(
		common.GenUUID(),
		dto.FirstName,
		dto.LastName,
		dto.Email,
		hashedPassword,
		salt,
		domain.RoleUser,
	)

	if err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, userEntity); err != nil {
		return err
	}

	return nil
}

type UserRepository interface {
	//Find(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, data *domain.User) error
	//Update(ctx context.Context, data *domain.User) error
	//Delete(ctx context.Context, data *domain.User) error
}
