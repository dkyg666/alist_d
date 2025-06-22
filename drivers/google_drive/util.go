package google_drive

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/vscodev/alist/v3/drivers/base"
	"github.com/vscodev/alist/v3/internal/driver"
	"github.com/vscodev/alist/v3/internal/model"
	"github.com/vscodev/alist/v3/pkg/http_range"
	"github.com/vscodev/alist/v3/pkg/utils"
)

func (d *GoogleDrive) refreshToken() error {
	url := d.OauthTokenURL
	if d.ClientID != "" && d.ClientSecret != "" {
		url = "https://oauth2.googleapis.com/token"
	}

	var resp base.TokenResp
	var e TokenError
	res, err := base.RestyClient.R().SetResult(&resp).SetError(&e).
		SetFormData(map[string]string{
			"client_id":     d.ClientID,
			"client_secret": d.ClientSecret,
			"refresh_token": d.RefreshToken,
			"grant_type":    "refresh_token",
		}).Post(url)
	if err != nil {
		return err
	}
	log.Debug(res.String())
	if e.Error != "" {
		return fmt.Errorf(e.Error)
	}
	d.AccessToken = resp.AccessToken
	return nil
}

func (d *GoogleDrive) request(url string, method string, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	req := base.RestyClient.R()
	req.SetHeader("Authorization", "Bearer "+d.AccessToken)
	req.SetQueryParam("includeItemsFromAllDrives", "true")
	req.SetQueryParam("supportsAllDrives", "true")
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
	var e Error
	req.SetError(&e)
	res, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}
	if e.Error.Code != 0 {
		if e.Error.Code == 401 {
			err = d.refreshToken()
			if err != nil {
				return nil, err
			}
			return d.request(url, method, callback, resp)
		}
		return nil, fmt.Errorf("%s: %v", e.Error.Message, e.Error.Errors)
	}
	return res.Body(), nil
}

func (d *GoogleDrive) getFiles(id string) ([]File, error) {
	pageToken := "first"
	res := make([]File, 0)
	for pageToken != "" {
		if pageToken == "first" {
			pageToken = ""
		}
		var resp Files
		orderBy := "folder,name,modifiedTime desc"
		if d.OrderBy != "" {
			orderBy = d.OrderBy + " " + d.OrderDirection
		}
		query := map[string]string{
			"orderBy":  orderBy,
			"fields":   "files(id,name,mimeType,size,modifiedTime,createdTime,thumbnailLink,shortcutDetails,md5Checksum,sha1Checksum,sha256Checksum),nextPageToken",
			"pageSize": "1000",
			"q":        fmt.Sprintf("'%s' in parents and trashed = false", id),
			// "includeItemsFromAllDrives": "true",
			// "supportsAllDrives":         "true",
			"pageToken": pageToken,
		}
		_, err := d.request("https://www.googleapis.com/drive/v3/files", http.MethodGet, func(req *resty.Request) {
			req.SetQueryParams(query)
		}, &resp)
		if err != nil {
			return nil, err
		}
		pageToken = resp.NextPageToken
		res = append(res, resp.Files...)
	}
	return res, nil
}

func (d *GoogleDrive) chunkUpload(ctx context.Context, stream model.FileStreamer, url string) error {
	var defaultChunkSize = d.ChunkSize * 1024 * 1024
	var offset int64 = 0
	for offset < stream.GetSize() {
		if utils.IsCanceled(ctx) {
			return ctx.Err()
		}
		chunkSize := stream.GetSize() - offset
		if chunkSize > defaultChunkSize {
			chunkSize = defaultChunkSize
		}
		reader, err := stream.RangeRead(http_range.Range{Start: offset, Length: chunkSize})
		if err != nil {
			return err
		}
		reader = driver.NewLimitedUploadStream(ctx, reader)
		_, err = d.request(url, http.MethodPut, func(req *resty.Request) {
			req.SetHeaders(map[string]string{
				"Content-Length": strconv.FormatInt(chunkSize, 10),
				"Content-Range":  fmt.Sprintf("bytes %d-%d/%d", offset, offset+chunkSize-1, stream.GetSize()),
			}).SetBody(reader).SetContext(ctx)
		}, nil)
		if err != nil {
			return err
		}
		offset += chunkSize
	}
	return nil
}
