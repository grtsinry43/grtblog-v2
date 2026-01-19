package content

import (
	"crypto/md5"
	"encoding/hex"
)

func ArticleContentHash(title string, leadIn *string, content string) string {
	return hashContentParts(title, stringOrEmpty(leadIn), content)
}

func MomentContentHash(title string, summary string, content string) string {
	return hashContentParts(title, summary, content)
}

func PageContentHash(title string, description *string, content string) string {
	return hashContentParts(title, stringOrEmpty(description), content)
}

func hashContentParts(parts ...string) string {
	hasher := md5.New()
	for i, part := range parts {
		if i > 0 {
			_, _ = hasher.Write([]byte{0})
		}
		_, _ = hasher.Write([]byte(part))
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func stringOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
