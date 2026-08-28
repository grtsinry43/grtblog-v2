package contentexport

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
)

func TestClassifier(t *testing.T) {
	uploadDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(uploadDir, "2025/04/23"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uploadDir, "2025/04/23/exists.png"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	c := &classifier{uploadDir: uploadDir, publicHost: "blog.example.com"}

	cases := []struct {
		name string
		ref  string
		kind imageKind
		sub  string
	}{
		{"相对站内", "/uploads/pictures/x.png", imageSelf, "pictures/x.png"},
		{"本站域名绝对URL", "https://blog.example.com/uploads/2025/04/23/exists.png", imageSelf, "2025/04/23/exists.png"},
		{"他域但磁盘存在(兜底)", "https://cdn.example.com/uploads/2025/04/23/exists.png", imageSelf, "2025/04/23/exists.png"},
		{"本站域名文件缺失仍是站内", "https://blog.example.com/uploads/2025/04/23/missing.png", imageSelf, "2025/04/23/missing.png"},
		{"第三方uploads且缺失->外链", "https://evil.com/uploads/a.png", imageExternal, ""},
		{"OSS外链图片", "https://dogeoss.example.com/a.png", imageExternal, ""},
		{"外链带查询串", "https://dogeoss.example.com/a.png?x=1", imageExternal, ""},
		{"外链非图片", "https://github.com/x/y", imageIgnored, ""},
		{"data URI", "data:image/png;base64,AAA", imageIgnored, ""},
		{"mailto", "mailto:a@b.c", imageIgnored, ""},
		{"站内非uploads路径", "/api/foo", imageIgnored, ""},
		{"空串", "", imageIgnored, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kind, sub := c.classify(tc.ref)
			if kind != tc.kind || sub != tc.sub {
				t.Fatalf("classify(%q) = (%v, %q), want (%v, %q)", tc.ref, kind, sub, tc.kind, tc.sub)
			}
		})
	}
}

func TestExtractImageRefsSkipsCode(t *testing.T) {
	md := strings.Join([]string{
		"para ![alt](/uploads/a.png) text",
		"",
		"```",
		"![nope](/uploads/code.png)",
		"```",
		"",
		"inline `/uploads/inline.png` code",
		"",
		`<img src="/uploads/html.png">`,
		"",
		`::: year-card cover="/uploads/directive.png"`,
		"",
		"句末标点 /uploads/b.png.",
	}, "\n")
	spans := extractImageRefs(md)
	refs := make(map[string]bool)
	for _, sp := range spans {
		refs[sp.ref] = true
	}
	for _, want := range []string{"/uploads/a.png", "/uploads/html.png", "/uploads/directive.png", "/uploads/b.png"} {
		if !refs[want] {
			t.Errorf("应提取 %q，实际: %v", want, refs)
		}
	}
	for _, nope := range []string{"/uploads/code.png", "/uploads/inline.png", "/uploads/b.png."} {
		if refs[nope] {
			t.Errorf("不应提取 %q", nope)
		}
	}
}

func newTestResolver(t *testing.T, publicHost string) (*Resolver, string) {
	t.Helper()
	uploadDir := t.TempDir()
	workdir := t.TempDir()
	return NewResolver(config.ExportConfig{}, uploadDir, workdir, publicHost), uploadDir
}

func TestResolverSelfCopyAndDedup(t *testing.T) {
	r, uploadDir := newTestResolver(t, "")
	if err := os.MkdirAll(filepath.Join(uploadDir, "pictures"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uploadDir, "pictures/x.png"), []byte("pngdata"), 0o600); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	r.See(ctx, "/uploads/pictures/x.png", "article:1")
	r.See(ctx, "/uploads/pictures/x.png", "article:2")
	if rel := r.FinalRel("/uploads/pictures/x.png"); rel != "uploads/pictures/x.png" {
		t.Fatalf("FinalRel = %q", rel)
	}
	workdir := r.workdir
	data, err := os.ReadFile(filepath.Join(workdir, "uploads", "pictures", "x.png"))
	if err != nil || string(data) != "pngdata" {
		t.Fatalf("复制的文件: %v %q", err, data)
	}
	self, ext, failed := r.Stats()
	if self != 1 || ext != 0 || len(failed) != 0 {
		t.Fatalf("stats = (%d, %d, %v)", self, ext, failed)
	}
}

func TestResolverMissingSelfFileKeepsOriginal(t *testing.T) {
	r, _ := newTestResolver(t, "")
	ctx := context.Background()
	r.See(ctx, "/uploads/pictures/gone.png", "article:1")
	r.See(ctx, "/uploads/pictures/gone.png", "article:2") // 不应重复记失败
	if rel := r.FinalRel("/uploads/pictures/gone.png"); rel != "" {
		t.Fatalf("缺失文件应返回空, got %q", rel)
	}
	_, _, failed := r.Stats()
	if len(failed) != 1 {
		t.Fatalf("应只记 1 次失败, got %v", failed)
	}
}

func TestResolverContainsPath(t *testing.T) {
	r, uploadDir := newTestResolver(t, "")
	if err := os.WriteFile(filepath.Join(uploadDir, "secret.png"), []byte("s"), 0o600); err != nil {
		t.Fatal(err)
	}
	r.See(context.Background(), "/uploads/../secret.png", "article:1")
	// ../secret.png 被规范化为站内根下的 secret.png 并可复制，但绝不能逃出 workdir/uploads
	if _, err := os.Stat(filepath.Join(r.workdir, "secret.png")); !os.IsNotExist(err) {
		t.Fatalf("文件不应出现在 uploads 之外: %v", err)
	}
	if rel := r.FinalRel("/uploads/../secret.png"); rel != "" && strings.Contains(rel, "..") {
		t.Fatalf("相对引用不应含 ..: %q", rel)
	}
}

func TestApplySpansDepth(t *testing.T) {
	r, uploadDir := newTestResolver(t, "")
	if err := os.MkdirAll(filepath.Join(uploadDir, "pictures"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uploadDir, "pictures/x.png"), []byte("d"), 0o600); err != nil {
		t.Fatal(err)
	}
	md := "![a](/uploads/pictures/x.png) end"
	spans := extractImageRefs(md)
	if len(spans) != 1 {
		t.Fatalf("spans = %v", spans)
	}
	r.See(context.Background(), spans[0].ref, "t")
	out := applySpans(md, spans, r, 3)
	want := "![a](../../../uploads/pictures/x.png) end"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestRelRefWithDepth(t *testing.T) {
	if got := relRefWithDepth(1, "uploads/x.png"); got != "../uploads/x.png" {
		t.Errorf("depth1 got %q", got)
	}
	if got := relRefWithDepth(6, "uploads/a/b.png"); got != strings.Repeat("../", 6)+"uploads/a/b.png" {
		t.Errorf("depth6 got %q", got)
	}
}

func TestExternalName(t *testing.T) {
	a := externalName("https://x.com/a.png")
	b := externalName("https://x.com/b.png")
	if a == b {
		t.Fatal("不同 URL 应产生不同文件名")
	}
	if !strings.HasSuffix(a, ".png") {
		t.Fatalf("应保留图片扩展名: %q", a)
	}
	if len(strings.TrimSuffix(a, ".png")) != 16 {
		t.Fatalf("哈希段应为 16 字符: %q", a)
	}
	if got := externalName("https://x.com/noext"); !strings.HasSuffix(got, ".bin") {
		t.Fatalf("无扩展名应回退 .bin: %q", got)
	}
}

func TestHostFromURL(t *testing.T) {
	if got := hostFromURL("https://blog.example.com/path"); got != "blog.example.com" {
		t.Errorf("got %q", got)
	}
	if got := hostFromURL("not a url"); got != "" {
		t.Errorf("got %q", got)
	}
}
