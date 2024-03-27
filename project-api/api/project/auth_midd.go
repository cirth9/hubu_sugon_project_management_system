package project

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	common "test.com/project-common"
	"test.com/project-common/errs"
)

var ignores = []string{
	"project/login/register",
	"project/login",
	"project/login/getCaptcha",
	"project/organization",
	"project/auth/apply"}

func Auth() func(*gin.Context) {
	return func(c *gin.Context) {
		result := &common.Result{}
		uri := c.Request.RequestURI
		for _, v := range ignores {
			if strings.Contains(uri, v) {
				c.Next()
				return
			}
		}
		//判断此uri是否在用户的授权列表中
		a := NewAuth()
		nodes, err := a.GetAuthNodes(c)
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort()
			return
		}

		for _, v := range nodes {
			if strings.Contains(uri, v) {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusOK, result.Fail(403, "无权限操作"))
		c.Abort()
		return
	}
}
