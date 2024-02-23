package usecase

import (
	"context"
	"my-app/common"
	"my-app/module/user/domain"
	"time"
)

type loginEmailPasswordUC struct {
	userRepo      UserQueryRepository
	sessionRepo   SessionCommandRepository
	tokenProvider TokenProvider
	hasher        Hasher
}

func NewLoginEmailPasswordUC(userRepo UserQueryRepository, sessionRepo SessionCommandRepository,
	tokenProvider TokenProvider, hasher Hasher) *loginEmailPasswordUC {
	return &loginEmailPasswordUC{userRepo: userRepo, sessionRepo: sessionRepo, tokenProvider: tokenProvider, hasher: hasher}
}

func (uc *loginEmailPasswordUC) LoginEmailPassword(ctx context.Context, dto EmailPasswordLoginDTO) (*TokenResponseDTO, error) {
	// 1. Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, dto.Email)

	if err != nil {
		return nil, err
	}

	// 2. Hash & compare password
	if ok := uc.hasher.CompareHashPassword(user.Password(), user.Salt(), dto.Password); !ok {
		return nil, domain.ErrInvalidEmailPassword
	}

	userId := user.Id()
	sessionId := common.GenUUID()

	// 3. Gen JWT
	accessToken, err := uc.tokenProvider.IssueToken(ctx, sessionId.String(), userId.String())

	if err != nil {
		return nil, err
	}

	// 4. Insert session into DB
	refreshToken, _ := uc.hasher.RandomStr(16)
	tokenExpAt := time.Now().UTC().Add(time.Second * time.Duration(uc.tokenProvider.TokenExpireInSeconds()))
	refreshExpAt := time.Now().UTC().Add(time.Second * time.Duration(uc.tokenProvider.RefreshExpireInSeconds()))

	session := domain.NewSession(sessionId, userId, refreshToken, tokenExpAt, refreshExpAt)

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	// 5. Return token response dto

	return &TokenResponseDTO{
		AccessToken:       accessToken,
		AccessTokenExpIn:  uc.tokenProvider.TokenExpireInSeconds(),
		RefreshToken:      refreshToken,
		RefreshTokenExpIn: uc.tokenProvider.RefreshExpireInSeconds(),
	}, nil
}
