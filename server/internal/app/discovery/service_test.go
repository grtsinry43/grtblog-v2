package discovery

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	domain "github.com/grtsinry43/grtblog-v2/server/internal/domain/discovery"
)

type testSettings map[string]string

func (s testSettings) WebsiteInfo(context.Context) (map[string]string, error) { return s, nil }
func (s testSettings) Timezone(context.Context) *time.Location {
	return time.FixedZone("Asia/Shanghai", 8*3600)
}

type testRepo struct {
	rows []domain.Record
	err  error
}

func (r *testRepo) List(context.Context) ([]domain.Record, error) { return r.rows, r.err }
func (r *testRepo) Document(_ context.Context, kind, slug string) (domain.Record, error) {
	for _, row := range r.rows {
		if row.Kind == kind && row.Slug == slug {
			return row, nil
		}
	}
	return domain.Record{}, domain.ErrNotFound
}

func TestCatalogCanonicalAndContentDates(t *testing.T) {
	created := time.Date(2026, 9, 5, 18, 0, 0, 0, time.UTC)
	repo := &testRepo{rows: []domain.Record{
		{Kind: "posts", Slug: "你好 & hello", Title: "A < B & C", ModifiedAt: created},
		{Kind: "moments", Slug: "note", CreatedAt: created, ModifiedAt: created},
		{Kind: "pages", Slug: "posts", Title: "builtin"},
		{Kind: "pages", Slug: "about", Title: "About"},
		{Kind: "posts", Slug: "../private", Title: "unsafe"},
	}}
	svc := NewService(repo, testSettings{"public_url": "https://blog.example/", "website_name": "Blog"})
	c, err := svc.Catalog(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	raw, err := c.Sitemap(0)
	if err != nil {
		t.Fatal(err)
	}
	var parsed urlSet
	if err := xml.Unmarshal(raw, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed.XMLName.Space != sitemapNamespace {
		t.Fatalf("namespace: %s", parsed.XMLName.Space)
	}
	seen := map[string]xmlURL{}
	for _, item := range parsed.URLs {
		if _, exists := seen[item.Loc]; exists {
			t.Fatalf("duplicate %s", item.Loc)
		}
		seen[item.Loc] = item
	}
	if _, ok := seen["https://blog.example/moments/2026/09/06/note/"]; !ok {
		t.Fatal("moment timezone mismatch")
	}
	if seen["https://blog.example/"].Lastmod != "" {
		t.Fatal("navigation has invented timestamp")
	}
	if seen["https://blog.example/moments/2026/09/06/note/"].Lastmod != "2026-09-05T18:00:00Z" {
		t.Fatal("incorrect content date")
	}
	if strings.Contains(string(raw), "private") {
		t.Fatal("invalid slug indexed")
	}
	if !strings.Contains(string(raw), "%E4%BD%A0%E5%A5%BD%20&amp;%20hello/") {
		t.Fatalf("encoding: %s", raw)
	}
}

func TestDocumentRequiresCanonicalDateAndPublicRecord(t *testing.T) {
	svc := NewService(&testRepo{rows: []domain.Record{{Kind: "moments", Slug: "note", Title: "笔记", Author: "作者", CreatedAt: time.Date(2026, 9, 5, 18, 0, 0, 0, time.UTC), Content: "正文"}, {Kind: "pages", Slug: "posts", Content: "should not leak"}}}, testSettings{"public_url": "https://blog.example"})
	for _, path := range []string{"/moments/2026/09/05/note/", "/posts/", "/missing/", "/internal/preview/foo/"} {
		if _, _, err := svc.Document(context.Background(), path); !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("%s: %v", path, err)
		}
	}
	body, canonical, err := svc.Document(context.Background(), "/moments/2026/09/06/note/")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, "作者：作者") || !strings.Contains(body, "正文") || canonical != "https://blog.example/moments/2026/09/06/note/" {
		t.Fatalf("%s %s", canonical, body)
	}
}

func TestLLMSConfigurationAndPagination(t *testing.T) {
	repo := &testRepo{}
	for i := 0; i < 101; i++ {
		repo.rows = append(repo.rows, domain.Record{Kind: "posts", Slug: fmt.Sprintf("post-%d", i), Title: fmt.Sprintf("Post [%d]", i), Summary: "summary\ncontinued"})
	}
	settings := testSettings{"public_url": "https://blog.example", "discovery.intro": "写作与生活", "discovery.author": "作者背景", "discovery.guidance": "请保留署名", "discovery.featured_paths": "/posts/post-0/\n/posts/withdrawn/"}
	svc := NewService(repo, settings)
	c, err := svc.Catalog(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	root, err := c.LLMS("", 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"写作与生活", "作者背景", "请保留署名", "## 精选阅读", "/posts/post-0/index.md"} {
		if !strings.Contains(root, text) {
			t.Fatalf("missing %q", text)
		}
	}
	if strings.Contains(root, "withdrawn") {
		t.Fatal("private/deleted featured URL leaked")
	}
	first, _ := c.LLMS("posts", 1)
	second, _ := c.LLMS("posts", 2)
	if strings.Count(first, "index.md") != 100 || strings.Count(second, "index.md") != 1 || !strings.Contains(first, "?page=2") || !strings.Contains(second, "?page=1") {
		t.Fatal("broken pagination")
	}
	for _, p := range []int{0, 3, int(^uint(0) >> 1)} {
		if _, err := c.LLMS("posts", p); !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("page %d: %v", p, err)
		}
	}
	if _, err := c.LLMS("bad", 1); !errors.Is(err, domain.ErrNotFound) {
		t.Fatal(err)
	}
	repo.rows = repo.rows[1:]
	c, _ = svc.Catalog(context.Background())
	root, _ = c.LLMS("", 1)
	if strings.Contains(root, "post-0") {
		t.Fatal("withdrawn featured entry cached")
	}
}

func TestSitemapShards(t *testing.T) {
	c := &Catalog{BaseURL: "https://blog.example"}
	for i := 0; i < SitemapPageSize+1; i++ {
		c.Entries = append(c.Entries, Entry{URL: fmt.Sprintf("https://blog.example/posts/%05d/", i)})
	}
	index, _ := c.Sitemap(0)
	if !strings.Contains(string(index), "<sitemapindex") || !strings.Contains(string(index), "/sitemaps/2.xml") {
		t.Fatal(string(index))
	}
	last, _ := c.Sitemap(2)
	var set urlSet
	if err := xml.Unmarshal(last, &set); err != nil || len(set.URLs) != 1 {
		t.Fatalf("%v %d", err, len(set.URLs))
	}
	if _, err := c.Sitemap(3); !errors.Is(err, domain.ErrNotFound) {
		t.Fatal(err)
	}
}

func TestCatalogFailsClosed(t *testing.T) {
	for _, raw := range []string{"", "javascript:alert(1)", "https://user:secret@blog.example", "https://blog.example/?secret=1"} {
		if _, err := NewService(&testRepo{}, testSettings{"public_url": raw}).Catalog(context.Background()); err == nil {
			t.Fatalf("accepted %q", raw)
		}
	}
	if _, err := NewService(&testRepo{err: errors.New("database unavailable")}, testSettings{"public_url": "https://blog.example"}).Catalog(context.Background()); err == nil {
		t.Fatal("returned success during outage")
	}
}

func TestCleanMarkdownPreservesCodeAndResolvesMedia(t *testing.T) {
	source := "::: callout title=\"注意\"\n正文\n:::\n\n::: link-card href=\"/about/\" title=\"关于\"\n:::\n\n![图](/uploads/a.png)\n[链接](../other/)\n[图][image]\n\n[image]: /uploads/ref.png\n\n`![代码](/uploads/code.png)`\n\n```md\n::: details summary=\"code\"\n![不改](/uploads/code.png)\n:::\n```\n\n    [不改](/code/)\n"
	got := CleanMarkdown(source, "https://blog.example/posts/hello/")
	for _, part := range []string{"> 注意", "正文", "[关于](<https://blog.example/about/>)", "![图](https://blog.example/uploads/a.png)", "[链接](https://blog.example/posts/other/)", "[image]: https://blog.example/uploads/ref.png", "`![代码](/uploads/code.png)`", "```md\n::: details summary=\"code\"\n![不改](/uploads/code.png)\n:::\n```", "    [不改](/code/)"} {
		if !strings.Contains(got, part) {
			t.Fatalf("missing %q in:\n%s", part, got)
		}
	}
}
