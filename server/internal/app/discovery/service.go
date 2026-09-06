package discovery

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	domain "github.com/grtsinry43/grtblog-v2/server/internal/domain/discovery"
)

type Settings interface {
	WebsiteInfo(context.Context) (map[string]string, error)
	Timezone(context.Context) *time.Location
}

type Service struct {
	repo     domain.Repository
	settings Settings
}

func NewService(repo domain.Repository, settings Settings) *Service { return &Service{repo, settings} }

type Entry struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Path        string `json:"path"`
	Kind        string `json:"kind"`
	Summary     string `json:"summary,omitempty"`
	Modified    string `json:"modified,omitempty"`
	MarkdownURL string `json:"markdownUrl,omitempty"`
}

type Catalog struct {
	SiteName string  `json:"siteName"`
	BaseURL  string  `json:"baseUrl"`
	Entries  []Entry `json:"entries"`
	info     map[string]string
}

// Reserved page slugs are actual application routes, not custom page bodies.
var reserved = map[string]bool{
	"posts": true, "moments": true, "pages": true, "categories": true, "columns": true,
	"tags": true, "timeline": true, "thinkings": true, "friends": true, "friends-timeline": true,
	"albums": true, "statistics": true, "sitemap": true, "sitemaps": true, "admin": true,
	"internal": true, "auth": true, "api": true, "feed": true, "og-image": true, "_app": true,
}

func validSlug(slug string) bool {
	return slug != "" && slug != "." && slug != ".." && !strings.ContainsAny(slug, "/\\?#\r\n\t")
}

func recordPath(row domain.Record, zone *time.Location) string {
	if !validSlug(row.Slug) {
		return ""
	}
	slug := url.PathEscape(row.Slug)
	switch row.Kind {
	case "posts", "categories", "columns", "albums":
		return "/" + row.Kind + "/" + slug + "/"
	case "photos":
		if row.ID <= 0 {
			return ""
		}
		return fmt.Sprintf("/albums/%s/photo/%d/", slug, row.ID)
	case "moments":
		return "/moments/" + row.CreatedAt.In(zone).Format("2006/01/02") + "/" + slug + "/"
	case "pages":
		if reserved[row.Slug] || strings.Contains(row.Slug, ".") {
			return ""
		}
		return "/" + slug + "/"
	}
	return ""
}

func baseURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u == nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("site.public_url must be an absolute HTTP(S) URL")
	}
	u.Path = ""
	u.RawPath = ""
	return strings.TrimRight(u.String(), "/"), nil
}

func (s *Service) Catalog(ctx context.Context) (*Catalog, error) {
	info, err := s.settings.WebsiteInfo(ctx)
	if err != nil {
		return nil, err
	}
	base, err := baseURL(info["public_url"])
	if err != nil {
		return nil, err
	}
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	c := &Catalog{SiteName: info["website_name"], BaseURL: base, info: info, Entries: make([]Entry, 0, len(rows)+10)}
	if c.SiteName == "" {
		c.SiteName = "个人博客"
	}
	for _, item := range []struct{ path, title string }{
		{"/", "首页"}, {"/posts/", "文章"}, {"/moments/", "手记"}, {"/thinkings/", "思考"},
		{"/timeline/", "时间线"}, {"/tags/", "标签"}, {"/friends/", "友链"}, {"/albums/", "相册"}, {"/sitemap/", "站点地图"},
	} {
		c.Entries = append(c.Entries, Entry{Title: item.title, Path: item.path, URL: base + item.path, Kind: "navigation"})
	}
	zone := s.settings.Timezone(ctx)
	seen := make(map[string]bool)
	for _, entry := range c.Entries {
		seen[entry.Path] = true
	}
	for _, row := range rows {
		path := recordPath(row, zone)
		if path == "" || seen[path] {
			continue
		}
		seen[path] = true
		entry := Entry{Title: row.Title, Path: path, URL: base + path, Kind: row.Kind, Summary: row.Summary}
		if !row.ModifiedAt.IsZero() {
			entry.Modified = row.ModifiedAt.UTC().Format(time.RFC3339)
		}
		if row.Kind == "posts" || row.Kind == "moments" || row.Kind == "pages" {
			entry.MarkdownURL = entry.URL + "index.md"
		}
		c.Entries = append(c.Entries, entry)
	}
	return c, nil
}

func (s *Service) Document(ctx context.Context, path string) (string, string, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	kind, slug := "", ""
	switch {
	case len(parts) == 2 && parts[0] == "posts":
		kind, slug = "posts", parts[1]
	case len(parts) == 5 && parts[0] == "moments":
		kind, slug = "moments", parts[4]
	case len(parts) == 1:
		kind, slug = "pages", parts[0]
	default:
		return "", "", domain.ErrNotFound
	}
	if !validSlug(slug) {
		return "", "", domain.ErrNotFound
	}
	row, err := s.repo.Document(ctx, kind, slug)
	if err != nil {
		return "", "", err
	}
	canonicalPath := recordPath(row, s.settings.Timezone(ctx))
	// Dates and reserved routes must agree with the human-readable page.
	decoded, _ := url.PathUnescape(canonicalPath)
	if canonicalPath == "" || decoded != "/"+strings.Trim(path, "/")+"/" {
		return "", "", domain.ErrNotFound
	}
	info, err := s.settings.WebsiteInfo(ctx)
	if err != nil {
		return "", "", err
	}
	base, err := baseURL(info["public_url"])
	if err != nil {
		return "", "", err
	}
	canonical := base + canonicalPath
	var out strings.Builder
	fmt.Fprintf(&out, "# %s\n\n原文：%s\n\n", inline(row.Title), canonical)
	if row.Author != "" {
		fmt.Fprintf(&out, "作者：%s\n\n", inline(row.Author))
	}
	fmt.Fprintf(&out, "发布时间：%s\n\n", row.CreatedAt.Format(time.RFC3339))
	if !row.ModifiedAt.IsZero() {
		fmt.Fprintf(&out, "内容更新时间：%s\n\n", row.ModifiedAt.Format(time.RFC3339))
	}
	out.WriteString(CleanMarkdown(row.Content, canonical))
	out.WriteString("\n")
	return out.String(), canonical, nil
}

// Robots depends only on settings, so catalog query failures do not prevent
// crawlers from retrieving the configured policy.
func (s *Service) Robots(ctx context.Context) (string, error) {
	info, err := s.settings.WebsiteInfo(ctx)
	if err != nil {
		return "", err
	}
	base, err := baseURL(info["public_url"])
	if err != nil {
		return "", err
	}
	policy := strings.TrimSpace(strings.ReplaceAll(info["discovery.robots"], "\r\n", "\n"))
	if policy == "" {
		policy = "User-agent: *\nDisallow: /admin"
	}
	return policy + "\n\nSitemap: " + base + "/sitemap.xml\n", nil
}

func inline(s string) string { return strings.Join(strings.Fields(s), " ") }
func markdownLabel(s string) string {
	return strings.NewReplacer("\\", "\\\\", "[", "\\[", "]", "\\]").Replace(inline(s))
}
func sortedEntries(entries []Entry) []Entry {
	result := append([]Entry(nil), entries...)
	sort.Slice(result, func(i, j int) bool { return result[i].URL < result[j].URL })
	return result
}
