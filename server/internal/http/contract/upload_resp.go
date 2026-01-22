package contract

import (
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/media"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

type UploadFileResp struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	PublicURL  string    `json:"publicUrl"`
	Type       string    `json:"type"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"createdAt"`
	Duplicated bool      `json:"duplicated"`
}

type UploadFileListResp struct {
	Items []UploadFileResp `json:"items"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

// UploadFileRespEnvelope 用于 swagger 展示。
type UploadFileRespEnvelope struct {
	Code   int            `json:"code"`
	BizErr string         `json:"bizErr"`
	Msg    string         `json:"msg"`
	Data   UploadFileResp `json:"data"`
	Meta   response.Meta  `json:"meta"`
}

// UploadFileListRespEnvelope 用于 swagger 展示。
type UploadFileListRespEnvelope struct {
	Code   int                `json:"code"`
	BizErr string             `json:"bizErr"`
	Msg    string             `json:"msg"`
	Data   UploadFileListResp `json:"data"`
	Meta   response.Meta      `json:"meta"`
}

type UploadRenameReq struct {
	Name string `json:"name"`
}

func ToUploadFileResp(file media.UploadFile, duplicated bool) UploadFileResp {
	publicURL := ""
	if file.Path != "" {
		publicURL = "/uploads" + file.Path
	}
	return UploadFileResp{
		ID:         file.ID,
		Name:       file.Name,
		Path:       file.Path,
		PublicURL:  publicURL,
		Type:       file.Type,
		Size:       file.Size,
		CreatedAt:  file.CreatedAt,
		Duplicated: duplicated,
	}
}
