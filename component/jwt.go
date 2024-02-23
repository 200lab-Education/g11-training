package component

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
)

const (
	defaultSecret                 = "very-important-please-change-it!" // in 32 bytes
	defaultExpireTokenInSeconds   = 60 * 60 * 24 * 7                   // 7d
	defaultExpireRefreshInSeconds = 60 * 60 * 24 * 14                  // 14d
)

var (
	ErrSecretKeyNotValid     = errors.New("secret key must be in 32 bytes")
	ErrTokenLifeTimeTooShort = errors.New("token life time too short")
)

type jwtx struct {
	id                     string
	secret                 string
	expireTokenInSeconds   int
	expireRefreshInSeconds int
}

func NewJWTProvider(secret string, expireTokenInSeconds, expireRefreshInSeconds int) *jwtx {
	return &jwtx{
		secret:                 secret,
		expireTokenInSeconds:   expireTokenInSeconds,
		expireRefreshInSeconds: expireRefreshInSeconds,
	}
}

func NewJWT(id string) *jwtx {
	return &jwtx{id: id}
}

func (j *jwtx) ID() string {
	return j.id
}

func (j *jwtx) InitFlags() {
	flag.StringVar(
		&j.secret,
		"jwt-secret",
		defaultSecret,
		"Secret key to sign JWT",
	)

	flag.IntVar(
		&j.expireTokenInSeconds,
		"jwt-exp-secs",
		defaultExpireTokenInSeconds,
		"Number of seconds access token will expired",
	)

	flag.IntVar(
		&j.expireRefreshInSeconds,
		"refresh-exp-secs",
		defaultExpireRefreshInSeconds,
		"Number of seconds refresh token will expired",
	)
}

//
//func (j *jwtx) Activate(_ sctx.ServiceContext) error {
//	if len(j.secret) < 32 {
//		return errors.WithStack(ErrSecretKeyNotValid)
//	}
//
//	if j.expireTokenInSeconds <= 60 {
//		return errors.WithStack(ErrTokenLifeTimeTooShort)
//	}
//
//	return nil
//}
//
//func (j *jwtx) Stop() error {
//	return nil
//}

func (j *jwtx) IssueToken(ctx context.Context, id, sub string) (token string, err error) {
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(j.expireTokenInSeconds))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        id,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSignedStr, err := t.SignedString([]byte(j.secret))

	if err != nil {
		return "", errors.WithStack(err)
	}

	return tokenSignedStr, nil
}

func (j *jwtx) TokenExpireInSeconds() int   { return j.expireTokenInSeconds }
func (j *jwtx) RefreshExpireInSeconds() int { return j.expireRefreshInSeconds }

func (j *jwtx) ParseToken(ctx context.Context, tokenString string) (claims *jwt.RegisteredClaims, err error) {
	var rc jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenString, &rc, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !token.Valid {
		return nil, errors.WithStack(err)
	}

	return &rc, nil
}
