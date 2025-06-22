package dropbox

import (
	"github.com/vscodev/alist/v3/internal/driver"
	"github.com/vscodev/alist/v3/internal/op"
)

const (
	DefaultClientID = "76lrwrklhdn1icb"
)

type Addition struct {
	RefreshToken string `json:"refresh_token" required:"true"`
	driver.RootPath

	OauthTokenURL string `json:"oauth_token_url"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`

	AccessToken     string
	RootNamespaceId string
}

var config = driver.Config{
	Name:              "Dropbox",
	LocalSort:         false,
	OnlyLocal:         false,
	OnlyProxy:         false,
	NoCache:           false,
	NoUpload:          false,
	NeedMs:            false,
	DefaultRoot:       "",
	NoOverwriteUpload: true,
}

func init() {
	op.RegisterDriver(func() driver.Driver {
		return &Dropbox{
			base:        "https://api.dropboxapi.com",
			contentBase: "https://content.dropboxapi.com",
		}
	})
}
