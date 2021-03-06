/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/11 18:11
 **/
package middleware

import (
	"github.com/coder2z/ndisk/pkg/recaptcha"
	R "github.com/coder2z/ndisk/pkg/response"
	"github.com/gin-gonic/gin"
)

func Recaptcha() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		return
		if !recaptcha.Verify(ctx.GetHeader("captcha")).Success {
			ctx.Abort()
			R.HandleCaptchaError(ctx)
			return
		}
		ctx.Next()
		return
	}
}
