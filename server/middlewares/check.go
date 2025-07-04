package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vscodev/alist/v3/internal/conf"
	"github.com/vscodev/alist/v3/pkg/utils"
	"github.com/vscodev/alist/v3/server/common"
)

func StoragesLoaded(c *gin.Context) {
	if conf.StoragesLoaded {
		c.Next()
	} else {
		if utils.SliceContains([]string{"", "/", "/favicon.ico"}, c.Request.URL.Path) {
			c.Next()
			return
		}
		paths := []string{"/assets", "/images", "/streamer", "/static"}
		for _, path := range paths {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}
		common.ErrorStrResp(c, "Loading storage, please wait", 500)
		c.Abort()
	}
}
