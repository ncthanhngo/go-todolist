package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"todolist/common"
	"todolist/component/tokenprovider"
	"todolist/module/user/model"

	"github.com/gin-gonic/gin"
)

type AuthenStore interface {
	FindUser(ctx context.Context, cond map[string]interface{}, moreInfo ...string) (*model.User, error)
}

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("Wrong authen header"),
		fmt.Sprintf("ErrWrongAuthHeader"))
}
func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader(nil)
	}
	return parts[1], nil
}

// RequiredAuth
// 1.Get Token from Header
// 2.Validate Token and Parse to payload
// 3.From the Token payload, we use user_id to find from DB
func RequireAuth(authStore AuthenStore, tokenProvider tokenprovider.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token, err := extractTokenFromHeaderString(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(401, err)
			return
		}

		payload, err := tokenProvider.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(401, err)
			return
		}
		user, err := authStore.FindUser(
			c.Request.Context(),
			map[string]interface{}{"id": payload.UserId()},
		)
		if err != nil {
			c.AbortWithStatusJSON(401, err)
			return
		}
		if user.Status == 0 {
			c.AbortWithStatusJSON(403,
				common.ErrNoPermission(errors.New("user banned or deleted")),
			)
			return
		}
		c.Set(common.CurrentUser, user)
		c.Next()
	}
}
