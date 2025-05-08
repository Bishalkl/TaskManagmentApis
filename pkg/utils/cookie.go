package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetRefreshTokenCookie(ctx *gin.Context, token string, maxAge int) {
	cookie := fmt.Sprintf(
		"refresh_token=%s; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=%d",
		token,
		maxAge,
	)
	ctx.Header("Set-Cookie", cookie)
}
