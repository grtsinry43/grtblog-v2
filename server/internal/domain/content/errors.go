package content

import "errors"

var ErrArticleNotFound = errors.New("article not found")
var ErrCategoryNotFound = errors.New("category not found")
var ErrTagNotFound = errors.New("tag not found")
var ErrArticleShortURLExists = errors.New("article short url exists")
