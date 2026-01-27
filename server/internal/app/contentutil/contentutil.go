package contentutil

// 把公共逻辑抽取了下，文章手记页面都能用欸

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/mozillazg/go-pinyin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type heading struct {
	level int
	text  string
}

var markdownParser = goldmark.New()
var tocAnchorSanitizer = regexp.MustCompile(`[^a-zA-Z0-9-]+`)
var shortURLSanitizer = regexp.MustCompile(`[^a-zA-Z0-9-]+`)

const (
	CommentAreaTypeArticle  = "article"
	CommentAreaTypeMoment   = "moment"
	CommentAreaTypePage     = "page"
	CommentAreaTypeThinking = "thinking"
)

func BuildCommentAreaName(areaType, title string) string {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return truncateRunes("评论区："+areaType, 255)
	}
	return truncateRunes("评论区："+areaType+"："+trimmed, 255)
}

func GenerateTOC(markdown string) []content.TOCNode {
	headings := extractHeadings(markdown)
	if len(headings) == 0 {
		return nil
	}

	var roots []*content.TOCNode
	type stackItem struct {
		level int
		node  *content.TOCNode
	}
	var stack []stackItem
	anchorCounts := make(map[string]int)
	for _, h := range headings {
		anchor := anchorFromHeading(h.text, anchorCounts)
		node := &content.TOCNode{
			Name:   h.text,
			Anchor: anchor,
		}

		for len(stack) > 0 && stack[len(stack)-1].level >= h.level {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			roots = append(roots, node)
		} else {
			parent := stack[len(stack)-1].node
			parent.Children = append(parent.Children, content.TOCNode{
				Name:   node.Name,
				Anchor: node.Anchor,
			})
			node = &parent.Children[len(parent.Children)-1]
		}
		stack = append(stack, stackItem{level: h.level, node: node})
	}

	result := make([]content.TOCNode, len(roots))
	for i, node := range roots {
		result[i] = *node
	}
	return result
}

func BuildSummary(summary, content string) string {
	if strings.TrimSpace(summary) != "" {
		return summary
	}
	return truncateRunes(content, 200)
}

func GenerateShortURLFromTitle(title string) string {
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

func GenerateRandomShortURL() string {
	bytes := make([]byte, 4)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func extractHeadings(markdown string) []heading {
	var headings []heading
	source := []byte(markdown)
	doc := markdownParser.Parser().Parse(text.NewReader(source))
	_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		headingNode, ok := node.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}
		headingText := strings.TrimSpace(extractHeadingText(headingNode, source))
		if headingText == "" {
			return ast.WalkContinue, nil
		}
		headings = append(headings, heading{
			level: headingNode.Level,
			text:  headingText,
		})
		return ast.WalkContinue, nil
	})
	return headings
}

func extractHeadingText(node *ast.Heading, source []byte) string {
	var builder strings.Builder
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch typed := n.(type) {
		case *ast.Text:
			value := string(typed.Segment.Value(source))
			value = strings.ReplaceAll(value, "\n", " ")
			builder.WriteString(value)
			if typed.SoftLineBreak() || typed.HardLineBreak() {
				builder.WriteByte(' ')
			}
		case *ast.String:
			value := string(typed.Value)
			value = strings.ReplaceAll(value, "\n", " ")
			builder.WriteString(value)
		}
		return ast.WalkContinue, nil
	})
	return builder.String()
}

func anchorFromHeading(text string, counts map[string]int) string {
	slug := slugifyHeading(text)
	if slug == "" {
		slug = "section"
	}
	seq := counts[slug]
	counts[slug] = seq + 1
	if seq == 0 {
		return slug
	}
	return fmt.Sprintf("%s-%d", slug, seq+1)
}

func slugifyHeading(input string) string {
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal

	var builder strings.Builder
	needsDash := false
	for _, r := range input {
		if r >= 0x4e00 && r <= 0x9fa5 {
			py := pinyin.Pinyin(string(r), args)
			if len(py) == 0 || len(py[0]) == 0 || py[0][0] == "" {
				continue
			}
			if needsDash && builder.Len() > 0 {
				builder.WriteByte('-')
			}
			builder.WriteString(py[0][0])
			needsDash = true
			continue
		}
		if unicode.IsSpace(r) || r == '-' {
			needsDash = true
			continue
		}
		if r > unicode.MaxASCII {
			continue
		}
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			if needsDash && builder.Len() > 0 {
				builder.WriteByte('-')
			}
			builder.WriteRune(unicode.ToLower(r))
			needsDash = false
		}
	}

	slug := tocAnchorSanitizer.ReplaceAllString(builder.String(), "-")
	slug = strings.Trim(slug, "-")
	slug = strings.ToLower(slug)
	return slug
}

func truncateRunes(input string, limit int) string {
	runes := []rune(input)
	if len(runes) <= limit {
		return input
	}
	return string(runes[:limit])
}
