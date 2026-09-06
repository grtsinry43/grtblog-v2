package discovery

import (
	"html"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/federation"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var markdownParser = goldmark.New()
var containerLine = regexp.MustCompile(`(?m)^ {0,3}(?:>[ \t]*)*:::[ \t]*([^\r\n]*)`)
var markdownDestination = regexp.MustCompile(`\]\([ \t]*(<[^>\n]+>|(?:\\.|[^\s()\\]|\([^()\n]*\))+)`)
var referenceDestination = regexp.MustCompile(`(?m)^ {0,3}\[[^\]\n]+\]:[ \t]*(<[^>\n]+>|[^\s]+)`)
var htmlDestination = regexp.MustCompile(`(?:src|href)=["']([^"']+)["']`)

type edit struct {
	start, end int
	value      string
}

// CleanMarkdown keeps Markdown source and code verbatim while removing our
// presentation-only container syntax and resolving relative media/link URLs.
func CleanMarkdown(source, canonical string) string {
	raw := []byte(source)
	doc := markdownParser.Parser().Parse(text.NewReader(raw))
	var code [][2]int
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindFencedCodeBlock, ast.KindCodeBlock:
			for i := 0; i < n.Lines().Len(); i++ {
				line := n.Lines().At(i)
				code = append(code, [2]int{line.Start, line.Stop})
			}
			return ast.WalkSkipChildren, nil
		case ast.KindCodeSpan:
			for child := n.FirstChild(); child != nil; child = child.NextSibling() {
				if t, ok := child.(*ast.Text); ok {
					code = append(code, [2]int{t.Segment.Start, t.Segment.Stop})
				}
			}
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	protected := func(start, end int) bool {
		for _, iv := range code {
			if start < iv[1] && end > iv[0] {
				return true
			}
		}
		return false
	}
	base, _ := url.Parse(canonical)
	absolute := func(raw string) string {
		u, err := url.Parse(raw)
		if err != nil || base == nil {
			return raw
		}
		if u.IsAbs() {
			return raw
		}
		return base.ResolveReference(u).String()
	}
	var edits []edit
	for _, m := range containerLine.FindAllStringSubmatchIndex(source, -1) {
		if protected(m[0], m[1]) {
			continue
		}
		value := exportComponent(source[m[2]:m[3]], absolute)
		prefix := source[m[0] : m[0]+strings.Index(source[m[0]:m[1]], ":::")]
		if strings.Contains(prefix, ">") {
			value = prefix + strings.ReplaceAll(value, "\n", "\n"+prefix)
		}
		edits = append(edits, edit{m[0], m[1], value})
	}
	for _, match := range federation.MentionPatternPublic().FindAllStringSubmatchIndex(source, -1) {
		if protected(match[0], match[1]) {
			continue
		}
		value := "@" + source[match[2]:match[3]] + "@" + source[match[4]:match[5]]
		edits = append(edits, edit{match[0], match[1], value})
	}
	for _, match := range expandedMention.FindAllStringSubmatchIndex(source, -1) {
		if protected(match[0], match[1]) {
			continue
		}
		edits = append(edits, edit{match[0], match[1], html.UnescapeString(source[match[2]:match[3]])})
	}
	for _, pattern := range []*regexp.Regexp{markdownDestination, referenceDestination, htmlDestination} {
		for _, m := range pattern.FindAllStringSubmatchIndex(source, -1) {
			start, end := m[2], m[3]
			if protected(m[0], m[1]) {
				continue
			}
			value := source[start:end]
			if strings.HasPrefix(value, "<") && strings.HasSuffix(value, ">") {
				start++
				end--
				value = source[start:end]
			}
			resolved := absolute(value)
			if resolved != value {
				edits = append(edits, edit{start, end, resolved})
			}
		}
	}
	sort.SliceStable(edits, func(i, j int) bool { return edits[i].start < edits[j].start })
	var b strings.Builder
	pos := 0
	for _, e := range edits {
		if e.start < pos {
			continue
		}
		b.WriteString(source[pos:e.start])
		b.WriteString(e.value)
		pos = e.end
	}
	b.WriteString(source[pos:])
	return b.String()
}
