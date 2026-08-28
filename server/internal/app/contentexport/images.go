package contentexport

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	fedinfra "github.com/grtsinry43/grtblog-v2/server/internal/infra/federation"
)

// externalBudgetBytes 是单次导出任务允许下载的外链图片总字节上限，
// 超出后剩余外链记入失败清单并保留原始链接。
const externalBudgetBytes int64 = 2 << 30

type imageKind int

const (
	imageIgnored imageKind = iota
	imageSelf
	imageExternal
)

// classifier 判定一个图片引用是站内、外链还是忽略。
//
// 规则（存在性兜底）：
//   - 以 /uploads/ 开头的相对引用 → 站内；
//   - 绝对 http(s) URL 且路径含 /uploads/ → 仅当 host 是站点 public_url 的 host，
//     或文件确实存在于 UploadDir 下，才判为站内（避免把第三方 /uploads/ 路径误判，
//     也兼容老文章里 https://<本站域名>/uploads/... 的绝对写法）；
//   - 其余带图片扩展名的绝对 URL → 外链；其他一律忽略。
type classifier struct {
	uploadDir  string
	publicHost string
}

func (c *classifier) classify(ref string) (imageKind, string) {
	trimmed := strings.TrimSpace(ref)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") ||
		strings.HasPrefix(trimmed, "data:") || strings.HasPrefix(trimmed, "mailto:") ||
		strings.HasPrefix(trimmed, "javascript:") || strings.HasPrefix(trimmed, "tel:") {
		return imageIgnored, ""
	}
	if !strings.Contains(trimmed, "://") {
		if sub, ok := uploadsSubpath(splitQuery(trimmed)); ok {
			return imageSelf, sub
		}
		return imageIgnored, ""
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return imageIgnored, ""
	}
	if sub, ok := uploadsSubpath(parsed.Path); ok {
		if (c.publicHost != "" && strings.EqualFold(parsed.Hostname(), c.publicHost)) || c.localExists(sub) {
			return imageSelf, sub
		}
	}
	if hasImageExt(parsed.Path) {
		return imageExternal, ""
	}
	return imageIgnored, ""
}

// resolve 把站内 subpath 映射到磁盘路径，并做路径容器校验（拒绝 ../ 逃逸）。
// 返回 (磁盘路径, 规范化子路径, ok)。
func (c *classifier) resolve(sub string) (string, string, bool) {
	cleaned := strings.TrimPrefix(path.Clean("/"+sub), "/")
	if cleaned == "" || cleaned == "." {
		return "", "", false
	}
	joined := filepath.Join(c.uploadDir, filepath.FromSlash(cleaned))
	rel, err := filepath.Rel(c.uploadDir, joined)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", "", false
	}
	return joined, filepath.ToSlash(rel), true
}

func (c *classifier) localExists(sub string) bool {
	disk, _, ok := c.resolve(sub)
	if !ok {
		return false
	}
	info, err := os.Stat(disk)
	return err == nil && info.Mode().IsRegular()
}

func uploadsSubpath(p string) (string, bool) {
	if !strings.HasPrefix(p, "/uploads/") {
		return "", false
	}
	sub := strings.TrimPrefix(p, "/uploads/")
	if sub == "" {
		return "", false
	}
	return sub, true
}

func splitQuery(ref string) string {
	if i := strings.IndexAny(ref, "?#"); i >= 0 {
		return ref[:i]
	}
	return ref
}

// imageExts 将扩展名规范化（.jpeg → .jpg）。
var imageExts = map[string]string{
	".png": ".png", ".jpg": ".jpg", ".jpeg": ".jpg", ".gif": ".gif",
	".webp": ".webp", ".avif": ".avif", ".svg": ".svg", ".bmp": ".bmp", ".ico": ".ico",
}

func hasImageExt(p string) bool {
	_, ok := imageExts[strings.ToLower(path.Ext(p))]
	return ok
}

func imageExt(p string) string {
	if ext, ok := imageExts[strings.ToLower(path.Ext(p))]; ok {
		return ext
	}
	return ".bin"
}

func hostFromURL(rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return ""
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return ""
	}
	return parsed.Hostname()
}

// ---- 引用提取（goldmark 代码遮罩 + 正则）----

var exportMarkdown = goldmark.New()

// imageURLRe 匹配 URL 形态的 token：http(s) 链接与 /uploads/ 相对路径。
// 分类器随后决定哪些真正参与打包（外链需带图片扩展名）。
var imageURLRe = regexp.MustCompile("https?://[^\\s\"'`<>)\\]]+|/uploads/[^\\s\"'`<>)\\]]+")

type rawSpan struct {
	start, end int
	ref        string
}

// extractImageRefs 返回 markdown 中所有候选图片引用的字节区间，
// 代码块 / 行内代码内的引用被排除，避免破坏代码示例。
func extractImageRefs(markdown string) []rawSpan {
	masked := codeIntervals([]byte(markdown))
	var spans []rawSpan
	for _, loc := range imageURLRe.FindAllStringIndex(markdown, -1) {
		start, end := loc[0], loc[1]
		if overlapsAny(masked, start, end) {
			continue
		}
		ref := strings.TrimRight(markdown[start:end], ".,;:!?")
		if ref == "" {
			continue
		}
		spans = append(spans, rawSpan{start: start, end: start + len(ref), ref: ref})
	}
	return spans
}

// codeIntervals 收集源码中所有代码区域的字节区间。
func codeIntervals(source []byte) [][2]int {
	doc := exportMarkdown.Parser().Parse(text.NewReader(source))
	var intervals [][2]int
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindFencedCodeBlock, ast.KindCodeBlock:
			lines := n.Lines()
			for i := 0; i < lines.Len(); i++ {
				seg := lines.At(i)
				intervals = append(intervals, [2]int{seg.Start, seg.Stop})
			}
			return ast.WalkSkipChildren, nil
		case ast.KindCodeSpan:
			for child := n.FirstChild(); child != nil; child = child.NextSibling() {
				if child.Kind() == ast.KindText {
					seg := child.(*ast.Text).Segment
					intervals = append(intervals, [2]int{seg.Start, seg.Stop})
				}
			}
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return intervals
}

func overlapsAny(intervals [][2]int, start, end int) bool {
	for _, iv := range intervals {
		if start < iv[1] && end > iv[0] {
			return true
		}
	}
	return false
}

// ---- 解析与打包 ----

// Resolver 负责把图片引用落到导出包的 uploads/ 目录：
// 站内图片从 UploadDir 复制；外链图片经 SSRF 安全客户端下载。
// 两种布局（structured / flatten）必须共享同一个 Resolver 实例，
// 以保证 uploads/ 全局去重。
type Resolver struct {
	cfg        config.ExportConfig
	classifier *classifier
	workdir    string
	client     *http.Client

	mu         sync.Mutex
	selfSeen   map[string]bool   // 站内 subpath -> 已尝试处理（无论成功与否，避免重复复制/重复记失败）
	selfRel    map[string]string // 站内 subpath -> "uploads/<cleaned>"（仅复制成功者）
	extSeen    map[string]string // 外链 URL -> "uploads/external/<name>"
	extOK      map[string]bool   // 外链 URL -> 下载成功
	failed     []FailedDownload
	pending    []pendingDownload
	downloaded int64 // atomic，已下载外链总字节
}

type pendingDownload struct {
	url          string
	diskPath     string
	referencedBy string
}

func NewResolver(cfg config.ExportConfig, uploadDir, workdir, publicHost string) *Resolver {
	timeout := cfg.ExternalTimeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &Resolver{
		cfg:        cfg,
		classifier: &classifier{uploadDir: uploadDir, publicHost: publicHost},
		workdir:    workdir,
		client:     fedinfra.NewSafeHTTPClient(timeout),
		selfSeen:   make(map[string]bool),
		selfRel:    make(map[string]string),
		extSeen:    make(map[string]string),
		extOK:      make(map[string]bool),
	}
}

// See 登记一个引用：站内立即复制，外链加入待下载队列。重复引用自动去重。
func (r *Resolver) See(ctx context.Context, ref, referencedBy string) {
	kind, sub := r.classifier.classify(ref)
	switch kind {
	case imageSelf:
		r.mu.Lock()
		seen := r.selfSeen[sub]
		r.selfSeen[sub] = true
		r.mu.Unlock()
		if seen {
			return
		}
		r.copySelf(sub, referencedBy)
	case imageExternal:
		r.mu.Lock()
		if _, seen := r.extSeen[ref]; !seen {
			name := externalName(ref)
			rel := "uploads/external/" + name
			r.extSeen[ref] = rel
			r.pending = append(r.pending, pendingDownload{
				url:          ref,
				diskPath:     filepath.Join(r.workdir, filepath.FromSlash(rel)),
				referencedBy: referencedBy,
			})
		}
		r.mu.Unlock()
	}
}

func (r *Resolver) copySelf(sub, referencedBy string) {
	src, cleaned, ok := r.classifier.resolve(sub)
	if !ok {
		r.recordFailed("/uploads/"+sub, referencedBy, errors.New("非法的站内图片路径"))
		return
	}
	dst := filepath.Join(r.workdir, "uploads", filepath.FromSlash(cleaned))
	if err := copyRegularFile(src, dst); err != nil {
		r.recordFailed("/uploads/"+sub, referencedBy, err)
		return
	}
	r.mu.Lock()
	r.selfRel[sub] = "uploads/" + cleaned
	r.mu.Unlock()
}

// DownloadAll 以受限并发下载全部待处理外链。失败项记入失败清单。
func (r *Resolver) DownloadAll(ctx context.Context) {
	r.mu.Lock()
	pending := r.pending
	r.pending = nil
	workers := r.cfg.ExternalWorkers
	r.mu.Unlock()
	if workers <= 0 {
		workers = 4
	}
	if len(pending) == 0 {
		return
	}
	if workers > len(pending) {
		workers = len(pending)
	}
	jobs := make(chan pendingDownload)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				r.downloadOne(ctx, job)
			}
		}()
	}
	for _, job := range pending {
		jobs <- job
	}
	close(jobs)
	wg.Wait()
}

func (r *Resolver) downloadOne(ctx context.Context, job pendingDownload) {
	max := r.cfg.MaxExternalBytes
	if max <= 0 {
		max = 25 << 20
	}
	err := func() error {
		if atomic.LoadInt64(&r.downloaded) > externalBudgetBytes {
			return errors.New("外链下载总量超出预算")
		}
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, job.url, nil)
		if reqErr != nil {
			return reqErr
		}
		req.Header.Set("User-Agent", "GrtBlog-ContentExport/1.0")
		resp, doErr := r.client.Do(req)
		if doErr != nil {
			return doErr
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP %d", resp.StatusCode)
		}
		if err := os.MkdirAll(filepath.Dir(job.diskPath), 0o700); err != nil {
			return err
		}
		tmp := job.diskPath + ".part"
		file, openErr := os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if openErr != nil {
			return openErr
		}
		failedWrite := true
		defer func() {
			_ = file.Close()
			if failedWrite {
				_ = os.Remove(tmp)
			}
		}()
		written, copyErr := io.Copy(file, io.LimitReader(resp.Body, max+1))
		if copyErr != nil {
			return copyErr
		}
		if written > max {
			return fmt.Errorf("图片超过 %d 字节上限", max)
		}
		if err := file.Close(); err != nil {
			return err
		}
		failedWrite = false
		if err := os.Rename(tmp, job.diskPath); err != nil {
			return err
		}
		atomic.AddInt64(&r.downloaded, written)
		return nil
	}()
	if err != nil {
		r.recordFailed(job.url, job.referencedBy, err)
		return
	}
	r.mu.Lock()
	r.extOK[job.url] = true
	r.mu.Unlock()
}

func (r *Resolver) recordFailed(ref, referencedBy string, err error) {
	r.mu.Lock()
	r.failed = append(r.failed, FailedDownload{URL: ref, ReferencedBy: referencedBy, Error: err.Error()})
	r.mu.Unlock()
}

// FinalRel 返回引用对应的包根相对路径（"uploads/..."）；
// 空串表示保留原始引用（忽略项、缺失的站内文件、下载失败的外链）。
func (r *Resolver) FinalRel(ref string) string {
	kind, sub := r.classifier.classify(ref)
	r.mu.Lock()
	defer r.mu.Unlock()
	switch kind {
	case imageSelf:
		return r.selfRel[sub] // 不存在则为 ""
	case imageExternal:
		if r.extOK[ref] {
			return r.extSeen[ref]
		}
	}
	return ""
}

// Stats 返回打包成功的站内/外链图片数与失败清单。
func (r *Resolver) Stats() (self, external int64, failed []FailedDownload) {
	r.mu.Lock()
	defer r.mu.Unlock()
	self = int64(len(r.selfRel))
	external = int64(len(r.extOK))
	failed = append([]FailedDownload(nil), r.failed...)
	return
}

// ---- 改写工具 ----

// relRefWithDepth 把包根相对路径转成从 depth 层目录出发的相对引用。
func relRefWithDepth(depth int, archiveRel string) string {
	return strings.Repeat("../", depth) + archiveRel
}

// applySpans 用解析结果改写 markdown 中的图片引用（按区间顺序拼接，无需倒序）。
func applySpans(source string, spans []rawSpan, r *Resolver, depth int) string {
	if len(spans) == 0 {
		return source
	}
	var b strings.Builder
	b.Grow(len(source))
	cursor := 0
	for _, sp := range spans {
		b.WriteString(source[cursor:sp.start])
		if rel := r.FinalRel(sp.ref); rel != "" {
			b.WriteString(relRefWithDepth(depth, rel))
		} else {
			b.WriteString(sp.ref)
		}
		cursor = sp.end
	}
	b.WriteString(source[cursor:])
	return b.String()
}

// rewriteValue 改写单个裸 URL 值（cover / moment image[]）。
func rewriteValue(value string, r *Resolver, depth int) string {
	if value == "" {
		return value
	}
	if rel := r.FinalRel(value); rel != "" {
		return relRefWithDepth(depth, rel)
	}
	return value
}

func externalName(rawURL string) string {
	h := sha256.Sum256([]byte(rawURL))
	ext := ".bin"
	if parsed, err := url.Parse(rawURL); err == nil {
		ext = imageExt(parsed.Path)
	}
	return hex.EncodeToString(h[:8]) + ext
}

// copyRegularFile 复制普通文件（拒绝符号链接），O_EXCL 防止重复写。
func copyRegularFile(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("不是普通文件: %s", src)
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	failed := true
	defer func() {
		_ = out.Close()
		if failed {
			_ = os.Remove(dst)
		}
	}()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	if err := out.Sync(); err != nil {
		return err
	}
	failed = false
	return nil
}
