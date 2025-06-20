package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/vscodev/alist/v3/internal/conf"
	"github.com/vscodev/alist/v3/internal/errs"
	"github.com/vscodev/alist/v3/internal/setting"
	"github.com/vscodev/alist/v3/server/common"
)

func SearchIndex(c *gin.Context) {
	mode := setting.GetStr(conf.SearchIndex)
	if mode == "none" {
		common.ErrorResp(c, errs.SearchNotAvailable, 500)
		c.Abort()
	} else {
		c.Next()
	}
}
