package media

import "errors"

var ErrUploadFileNotFound = errors.New("上传文件不存在")
var ErrInvalidUploadType = errors.New("无效的上传类型")
