package contentexport

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
)

const (
	structuredRoot = "site-root"
	flattenRoot    = "site-flatten"
)

// preparedItem 是完成内容映射与引用登记后的待写入项。
// meta 指向原始（未改写）的 contract 详情结构体，写入时按目标布局的
// 目录深度对 cover / image[] / extInfo 做相对化改写，再剥离 content。
type preparedItem struct {
	kind        string // article | moment | page | thinking
	id          int64
	routePath   string
	sourcePath  string
	structDir   []string // structured 布局下相对 site-root 的目录段
	flattenName string   // flatten 布局下的文件名
	content     string
	spans       []rawSpan
	meta        any               // *contract.ArticleResp | *MomentResp | *PageResp | *ThinkingResp
	extInfo     *contract.JSONRaw // 与 meta 内 ExtInfo 同一引用
}

// metaJSON 生成 meta.json / flatten 头部的 MetaNode JSON。
// depth 为文件所在目录相对归档根的深度。
func (it *preparedItem) metaJSON(exportedAt time.Time, r *Resolver, depth int) ([]byte, error) {
	var metadata any
	switch v := it.meta.(type) {
	case *contract.ArticleResp:
		cp := *v
		if v.Cover != nil && *v.Cover != "" {
			rewritten := rewriteValue(*v.Cover, r, depth)
			cp.Cover = &rewritten
		}
		cp.ExtInfo = rewriteExtInfo(it.extInfo, r, depth)
		metadata = articleMeta{ArticleResp: &cp}
	case *contract.MomentResp:
		cp := *v
		if len(v.Image) > 0 {
			images := make([]string, len(v.Image))
			for i, img := range v.Image {
				images[i] = rewriteValue(img, r, depth)
			}
			cp.Image = images
		}
		cp.ExtInfo = rewriteExtInfo(it.extInfo, r, depth)
		metadata = momentMeta{MomentResp: &cp}
	case *contract.PageResp:
		cp := *v
		cp.ExtInfo = rewriteExtInfo(it.extInfo, r, depth)
		metadata = pageMeta{PageResp: &cp}
	case *contract.ThinkingResp:
		cp := *v
		metadata = thinkingMeta{ThinkingResp: &cp}
	default:
		return nil, fmt.Errorf("未知的导出项类型: %s", it.kind)
	}
	node := MetaNode{
		Kind:       it.kind,
		ID:         it.id,
		RoutePath:  it.routePath,
		SourcePath: it.sourcePath,
		ExportedAt: exportedAt,
		Metadata:   metadata,
	}
	return marshalIndent(node)
}

// WriteStructured 写入 site-root/ 下的 content.md + meta.json 结构。
func WriteStructured(items []*preparedItem, workdir string, r *Resolver, exportedAt time.Time) error {
	for _, it := range items {
		dirParts := append([]string{workdir, structuredRoot}, it.structDir...)
		dir := filepath.Join(dirParts...)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
		depth := len(dirParts) - 1 // workdir 本身是归档根
		content := applySpans(it.content, it.spans, r, depth)
		if err := os.WriteFile(filepath.Join(dir, "content.md"), []byte(content), 0o600); err != nil {
			return err
		}
		metaBytes, err := it.metaJSON(exportedAt, r, depth)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, "meta.json"), append(metaBytes, '\n'), 0o600); err != nil {
			return err
		}
	}
	return nil
}

// WriteFlatten 写入 site-flatten/ 下的单文件（---meta / ---content）。
func WriteFlatten(items []*preparedItem, workdir string, r *Resolver, exportedAt time.Time) error {
	dir := filepath.Join(workdir, flattenRoot)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	depth := 1 // site-flatten/<file>.md
	for _, it := range items {
		content := applySpans(it.content, it.spans, r, depth)
		metaBytes, err := it.metaJSON(exportedAt, r, depth)
		if err != nil {
			return err
		}
		doc := strings.Join([]string{"---meta", string(metaBytes), "---content", content, ""}, "\n")
		if err := os.WriteFile(filepath.Join(dir, it.flattenName), []byte(doc), 0o600); err != nil {
			return err
		}
	}
	return nil
}

type taxonomyEntry struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ShortURL string `json:"shortUrl,omitempty"`
}

// WriteTaxonomy 写入归档根下的 taxonomy.json（全部 tag/column/category 字典）。
func WriteTaxonomy(snap *Snapshot, workdir string) error {
	categories := make([]taxonomyEntry, 0, len(snap.Categories))
	for _, c := range snap.Categories {
		entry := taxonomyEntry{ID: c.ID, Name: c.Name}
		if c.ShortURL != nil {
			entry.ShortURL = *c.ShortURL
		}
		categories = append(categories, entry)
	}
	columns := make([]taxonomyEntry, 0, len(snap.Columns))
	for _, c := range snap.Columns {
		entry := taxonomyEntry{ID: c.ID, Name: c.Name}
		if c.ShortURL != nil {
			entry.ShortURL = *c.ShortURL
		}
		columns = append(columns, entry)
	}
	tags := make([]taxonomyEntry, 0, len(snap.Tags))
	for _, t := range snap.Tags {
		tags = append(tags, taxonomyEntry{ID: t.ID, Name: t.Name})
	}
	doc := struct {
		Categories []taxonomyEntry `json:"categories"`
		Columns    []taxonomyEntry `json:"columns"`
		Tags       []taxonomyEntry `json:"tags"`
	}{Categories: categories, Columns: columns, Tags: tags}
	raw, err := marshalIndent(doc)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(workdir, "taxonomy.json"), append(raw, '\n'), 0o600)
}

// ---- 命名工具（与旧 node 导出脚本对齐）----

var safeSlugRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

func slugFor(shortURL, kind string, id int64) string {
	s := strings.TrimSpace(shortURL)
	if s == "" {
		return fmt.Sprintf("%s-%d", kind, id)
	}
	if safeSlugRe.MatchString(s) {
		return s
	}
	return url.PathEscape(s)
}

var (
	flattenUnsafeRe = regexp.MustCompile(`[\\/:*?"<>|]`)
	flattenSpaceRe  = regexp.MustCompile(`\s+`)
)

func safeFileSegment(value, fallback string) string {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return fallback
	}
	s := flattenUnsafeRe.ReplaceAllString(raw, "-")
	s = flattenSpaceRe.ReplaceAllString(s, " ")
	s = strings.Trim(s, ".")
	if runes := []rune(s); len(runes) > 120 {
		s = string(runes[:120])
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return fallback
	}
	return s
}

func flattenFilename(kind, title string, id int64) string {
	fallback := fmt.Sprintf("%s-%d", kind, id)
	return fmt.Sprintf("%s__%s__%d.md", kind, safeFileSegment(title, fallback), id)
}

func dateParts(t time.Time, tz *time.Location) (string, string, string) {
	tt := t.In(tz)
	return fmt.Sprintf("%04d", tt.Year()), fmt.Sprintf("%02d", int(tt.Month())), fmt.Sprintf("%02d", tt.Day())
}

func firstLine(content string) string {
	if i := strings.IndexByte(content, '\n'); i >= 0 {
		return strings.TrimSpace(content[:i])
	}
	return strings.TrimSpace(content)
}
