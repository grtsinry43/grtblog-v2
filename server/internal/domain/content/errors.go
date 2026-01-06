package content

import "errors"

var ErrArticleNotFound = errors.New("文章不存在")
var ErrCategoryNotFound = errors.New("分类不存在")
var ErrColumnNotFound = errors.New("手记分区不存在")
var ErrTagNotFound = errors.New("标签不存在")
var ErrArticleShortURLExists = errors.New("文章短链接已存在")
