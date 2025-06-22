package yandex_disk

import (
	"github.com/vscodev/alist/v3/internal/driver"
	"github.com/vscodev/alist/v3/internal/op"
)

type Addition struct {
	RefreshToken   string `json:"refresh_token" required:"true"`
	OrderBy        string `json:"order_by" type:"select" options:"name,path,created,modified,size" default:"name"`
	OrderDirection string `json:"order_direction" type:"select" options:"asc,desc" default:"asc"`
	driver.RootPath
	OauthTokenURL string `json:"oauth_token_url"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
}

var config = driver.Config{
	Name:        "YandexDisk",
	DefaultRoot: "/",
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &YandexDisk{}
	})
}
