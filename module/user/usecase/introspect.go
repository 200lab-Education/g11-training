package usecase

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"my-app/common"
)

type TokenParser interface {
	ParseToken(ctx context.Context, tokenString string) (claims *jwt.RegisteredClaims, err error)
}

type introspectUC struct {
	userQueryRepo    UserQueryRepository
	sessionQueryRepo SessionQueryRepository
	tokenParser      TokenParser
}

func NewIntrospectUC(userQueryRepo UserQueryRepository, sessionQueryRepo SessionQueryRepository, tokenParser TokenParser) *introspectUC {
	return &introspectUC{userQueryRepo: userQueryRepo, sessionQueryRepo: sessionQueryRepo, tokenParser: tokenParser}
}

func (uc *introspectUC) IntrospectToken(ctx context.Context, accessToken string) (common.Requester, error) {
	claims, err := uc.tokenParser.ParseToken(ctx, accessToken)

	if err != nil {
		return nil, err
	}

	userId := uuid.MustParse(claims.Subject)
	sessionId := uuid.MustParse(claims.ID)

	if _, err := uc.sessionQueryRepo.Find(ctx, sessionId); err != nil {
		return nil, err
	}

	user, err := uc.userQueryRepo.Find(ctx, userId)

	if err != nil {
		return nil, err
	}

	if user.Status() == "banned" {
		return nil, errors.New("user has been banned")
	}

	return common.NewRequester(userId, sessionId, user.FirstName(), user.LastName(), user.Role().String(), user.Status()), nil
}
