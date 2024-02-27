package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"my-app/common"
	"strings"
)

type AuthClient interface {
	IntrospectToken(ctx context.Context, accessToken string) (common.Requester, error)
}

func RequireAuth(ac AuthClient) func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))

		if err != nil {
			common.WriteErrorResponse(c, err)
			c.Abort()
			return
		}

		requester, err := ac.IntrospectToken(c.Request.Context(), token)

		if err != nil {
			common.WriteErrorResponse(c, err)
			c.Abort()
			return
		}

		c.Set(common.KeyRequester, requester)

		c.Next()
	}
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	//"Authorization" : "Bearer {token}"

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("missing access token")
	}

	return parts[1], nil
}
