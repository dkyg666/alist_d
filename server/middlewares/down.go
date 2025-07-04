package middlewares

import (
	"strings"

	"github.com/vscodev/alist/v3/internal/conf"
	"github.com/vscodev/alist/v3/internal/setting"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/vscodev/alist/v3/internal/errs"
	"github.com/vscodev/alist/v3/internal/model"
	"github.com/vscodev/alist/v3/internal/op"
	"github.com/vscodev/alist/v3/pkg/utils"
	"github.com/vscodev/alist/v3/server/common"
)

func Down(verifyFunc func(string, string) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		rawPath := parsePath(c.Param("path"))
		c.Set("path", rawPath)
		meta, err := op.GetNearestMeta(rawPath)
		if err != nil {
			if !errors.Is(errors.Cause(err), errs.MetaNotFound) {
				common.ErrorResp(c, err, 500, true)
				return
			}
		}
		c.Set("meta", meta)
		// verify sign
		if needSign(meta, rawPath) {
			s := c.Query("sign")
			err = verifyFunc(rawPath, strings.TrimSuffix(s, "/"))
			if err != nil {
				common.ErrorResp(c, err, 401)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// TODO: implement
// path maybe contains # ? etc.
func parsePath(path string) string {
	return utils.FixAndCleanPath(path)
}

func needSign(meta *model.Meta, path string) bool {
	if setting.GetBool(conf.SignAll) {
		return true
	}
	if common.IsStorageSignEnabled(path) {
		return true
	}
	if meta == nil || meta.Password == "" {
		return false
	}
	if !meta.PSub && path != meta.Path {
		return false
	}
	return true
}
