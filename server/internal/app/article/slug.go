package article

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

var shortURLSanitizer = regexp.MustCompile(`[^a-zA-Z0-9-]+`)

func generateShortURLFromTitle(title string) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal

	var builder strings.Builder
	wordCount := 0
	for _, r := range title {
		if r >= 0x4e00 && r <= 0x9fa5 {
			py := pinyin.Pinyin(string(r), args)
			if len(py) == 0 || len(py[0]) == 0 || py[0][0] == "" {
				continue
			}
			if wordCount > 0 {
				builder.WriteByte('-')
			}
			builder.WriteString(py[0][0])
			wordCount++
		} else if unicode.IsSpace(r) {
			if wordCount > 0 {
				builder.WriteByte('-')
			}
			wordCount++
		} else {
			builder.WriteRune(r)
		}

		if wordCount >= 10 {
			break
		}
	}

	slug := shortURLSanitizer.ReplaceAllString(builder.String(), "")
	return strings.ToLower(slug)
}

func generateRandomShortURL() string {
	bytes := make([]byte, 4)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
