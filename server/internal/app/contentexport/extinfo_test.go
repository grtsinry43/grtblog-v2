package contentexport

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
)

func TestRewriteExtInfo(t *testing.T) {
	uploadDir := t.TempDir()
	workdir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(uploadDir, "pictures"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uploadDir, "pictures/x.png"), []byte("d"), 0o600); err != nil {
		t.Fatal(err)
	}
	r := NewResolver(config.ExportConfig{}, uploadDir, workdir, "")

	raw := contract.JSONRaw(`{"cover":"/uploads/pictures/x.png","nested":{"list":["/uploads/pictures/x.png","keep"],"key":"/uploads/pictures/x.png"},"num":1}`)
	for _, ref := range extInfoRefs(&raw) {
		r.See(context.Background(), ref, "article:1")
	}
	out := rewriteExtInfo(&raw, r, 2)
	if out == nil {
		t.Fatal("out 不应为 nil")
	}
	var v map[string]any
	if err := json.Unmarshal([]byte(*out), &v); err != nil {
		t.Fatal(err)
	}
	if v["cover"] != "../../uploads/pictures/x.png" {
		t.Fatalf("cover = %v", v["cover"])
	}
	nested := v["nested"].(map[string]any)
	list := nested["list"].([]any)
	if list[0] != "../../uploads/pictures/x.png" {
		t.Fatalf("list[0] = %v", list[0])
	}
	if list[1] != "keep" {
		t.Fatalf("非图片字符串不应改写: %v", list[1])
	}
	if nested["key"] != "../../uploads/pictures/x.png" {
		t.Fatalf("嵌套值 = %v", nested["key"])
	}
	if v["num"] != float64(1) {
		t.Fatalf("数值不应变: %v", v["num"])
	}
}

func TestRewriteExtInfoPassthrough(t *testing.T) {
	r := NewResolver(config.ExportConfig{}, t.TempDir(), t.TempDir(), "")

	invalid := contract.JSONRaw(`{invalid`)
	if out := rewriteExtInfo(&invalid, r, 1); string(*out) != string(invalid) {
		t.Fatal("非法 JSON 应原样返回")
	}
	if out := rewriteExtInfo(nil, r, 1); out != nil {
		t.Fatal("nil 应原样返回")
	}
}

func TestExtInfoRefs(t *testing.T) {
	raw := contract.JSONRaw(`{"images":[{"id":"https://blog.example.com/uploads/a.png"}],"plain":"text"}`)
	refs := extInfoRefs(&raw)
	if len(refs) != 1 || refs[0] != "https://blog.example.com/uploads/a.png" {
		t.Fatalf("refs = %v", refs)
	}
	if got := extInfoRefs(nil); len(got) != 0 {
		t.Fatalf("nil -> %v", got)
	}
}
