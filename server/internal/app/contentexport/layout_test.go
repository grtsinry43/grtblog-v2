package contentexport

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
)

func TestSlugFor(t *testing.T) {
	if got := slugFor("2025-summary", "article", 32); got != "2025-summary" {
		t.Errorf("got %q", got)
	}
	if got := slugFor("", "article", 32); got != "article-32" {
		t.Errorf("got %q", got)
	}
	if got := slugFor("奇怪 标题", "article", 5); strings.Contains(got, " ") || got == "" {
		t.Errorf("非 ASCII slug 应转义: %q", got)
	}
}

func TestFlattenFilename(t *testing.T) {
	name := flattenFilename("article", `a/b:c*d?e"f<g>h|i\j`, 7)
	if !strings.HasPrefix(name, "article__") || !strings.HasSuffix(name, "__7.md") {
		t.Fatalf("got %q", name)
	}
	if strings.ContainsAny(name, `\/:*?"<>|`) {
		t.Fatalf("文件名含非法字符: %q", name)
	}
	if got := flattenFilename("thinking", "", 3); got != "thinking__thinking-3__3.md" {
		t.Fatalf("空标题回退: %q", got)
	}
}

func TestDateParts(t *testing.T) {
	tz := time.FixedZone("CST", 8*3600)
	// UTC 2025-12-31 20:00 -> CST 2026-01-01 04:00，日期目录必须用站点时区
	ts := time.Date(2025, 12, 31, 20, 0, 0, 0, time.UTC)
	y, m, d := dateParts(ts, tz)
	if y != "2026" || m != "01" || d != "01" {
		t.Fatalf("got %s/%s/%s", y, m, d)
	}
}

func TestWriteFlattenRoundTrip(t *testing.T) {
	workdir := t.TempDir()
	r := NewResolver(config.ExportConfig{}, t.TempDir(), workdir, "")
	exportedAt := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)
	resp := &contract.ThinkingResp{ID: 1, Content: "body", CreatedAt: exportedAt, UpdatedAt: exportedAt}
	item := &preparedItem{
		kind:        "thinking",
		id:          1,
		routePath:   "/thinkings#thinking-1",
		sourcePath:  "/thinkings/1",
		flattenName: "thinking__x__1.md",
		content:     "body",
		meta:        resp,
	}
	if err := WriteFlatten([]*preparedItem{item}, workdir, r, exportedAt); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(workdir, flattenRoot, "thinking__x__1.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(raw)
	metaIdx := strings.Index(doc, "---meta\n")
	contentIdx := strings.Index(doc, "---content\n")
	if metaIdx != 0 || contentIdx < 0 {
		t.Fatalf("文档结构错误:\n%s", doc)
	}
	var node MetaNode
	if err := json.Unmarshal([]byte(doc[len("---meta\n"):contentIdx]), &node); err != nil {
		t.Fatalf("meta 解析失败: %v", err)
	}
	if node.Kind != "thinking" || node.ID != 1 {
		t.Fatalf("node = %+v", node)
	}
	metaRaw, _ := json.Marshal(node.Metadata)
	if strings.Contains(string(metaRaw), `"content"`) {
		t.Fatalf("meta 不应含 content 字段: %s", metaRaw)
	}
	if body := doc[contentIdx+len("---content\n"):]; !strings.HasPrefix(body, "body") {
		t.Fatalf("正文错误: %q", body)
	}
}

func TestWriteStructuredImageRewrite(t *testing.T) {
	uploadDir := t.TempDir()
	workdir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(uploadDir, "pictures"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uploadDir, "pictures/x.png"), []byte("d"), 0o600); err != nil {
		t.Fatal(err)
	}
	r := NewResolver(config.ExportConfig{}, uploadDir, workdir, "")
	exportedAt := time.Date(2026, 7, 20, 0, 0, 0, 0, time.UTC)

	cover := "/uploads/pictures/x.png"
	resp := &contract.ArticleResp{ID: 9, Title: "t", Content: "![](/uploads/pictures/x.png)", Cover: &cover, ShortURL: "demo"}
	item := &preparedItem{
		kind:        "article",
		id:          9,
		routePath:   "/posts/demo",
		sourcePath:  "/admin/articles/9",
		structDir:   []string{"posts", "demo"},
		flattenName: "article__t__9.md",
		content:     resp.Content,
		spans:       extractImageRefs(resp.Content),
		meta:        resp,
		extInfo:     resp.ExtInfo,
	}
	for _, sp := range item.spans {
		r.See(context.Background(), sp.ref, "article:9")
	}
	r.See(context.Background(), cover, "article:9")

	if err := WriteStructured([]*preparedItem{item}, workdir, r, exportedAt); err != nil {
		t.Fatal(err)
	}
	mdRaw, err := os.ReadFile(filepath.Join(workdir, structuredRoot, "posts", "demo", "content.md"))
	if err != nil {
		t.Fatal(err)
	}
	// depth = site-root(1) + posts(2) + demo(3)
	if !strings.Contains(string(mdRaw), "![](../../../uploads/pictures/x.png)") {
		t.Fatalf("正文引用改写错误: %s", mdRaw)
	}
	metaRaw, err := os.ReadFile(filepath.Join(workdir, structuredRoot, "posts", "demo", "meta.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(metaRaw), `"cover": "../../../uploads/pictures/x.png"`) {
		t.Fatalf("cover 改写错误: %s", metaRaw)
	}
	if strings.Contains(string(metaRaw), `"content"`) {
		t.Fatalf("meta.json 不应含 content: %s", metaRaw)
	}
	// 图片确实被复制进包
	if _, err := os.Stat(filepath.Join(workdir, "uploads", "pictures", "x.png")); err != nil {
		t.Fatalf("uploads 未打包: %v", err)
	}
}

func TestWriteTaxonomyEmptyArrays(t *testing.T) {
	workdir := t.TempDir()
	snap := &Snapshot{}
	if err := WriteTaxonomy(snap, workdir); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(workdir, "taxonomy.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"categories": []`, `"columns": []`, `"tags": []`} {
		if !strings.Contains(string(raw), key) {
			t.Fatalf("空站应输出空数组而非 null: %s", raw)
		}
	}
}
