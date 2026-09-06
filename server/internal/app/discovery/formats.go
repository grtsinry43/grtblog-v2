package discovery

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	domain "github.com/grtsinry43/grtblog-v2/server/internal/domain/discovery"
)

const SitemapPageSize = 10000
const DirectoryPageSize = 100

type xmlURL struct {
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod,omitempty"`
}
type urlSet struct {
	XMLName   xml.Name `xml:"urlset"`
	Namespace string   `xml:"xmlns,attr"`
	URLs      []xmlURL `xml:"url"`
}
type sitemapIndex struct {
	XMLName   xml.Name `xml:"sitemapindex"`
	Namespace string   `xml:"xmlns,attr"`
	Maps      []xmlURL `xml:"sitemap"`
}

const sitemapNamespace = "http://www.sitemaps.org/schemas/sitemap/0.9"

// Sitemap returns an index only when sharding is needed; page=0 is the root.
func (c *Catalog) Sitemap(page int) ([]byte, error) {
	entries := sortedEntries(c.Entries)
	if page == 0 && len(entries) > SitemapPageSize {
		index := sitemapIndex{Namespace: sitemapNamespace}
		for n := 1; n <= (len(entries)+SitemapPageSize-1)/SitemapPageSize; n++ {
			index.Maps = append(index.Maps, xmlURL{Loc: fmt.Sprintf("%s/sitemaps/%d.xml", c.BaseURL, n)})
		}
		body, err := xml.MarshalIndent(index, "", "  ")
		return append([]byte(xml.Header), body...), err
	}
	if page < 0 {
		return nil, domain.ErrNotFound
	}
	if page > 0 {
		if page > (len(entries)+SitemapPageSize-1)/SitemapPageSize {
			return nil, domain.ErrNotFound
		}
		start := (page - 1) * SitemapPageSize
		entries = entries[start:min(start+SitemapPageSize, len(entries))]
	}
	set := urlSet{Namespace: sitemapNamespace}
	for _, entry := range entries {
		set.URLs = append(set.URLs, xmlURL{entry.URL, entry.Modified})
	}
	body, err := xml.MarshalIndent(set, "", "  ")
	return append([]byte(xml.Header), body...), err
}

func (c *Catalog) LLMS(scope string, page int) (string, error) {
	if scope != "" && scope != "posts" && scope != "moments" && scope != "pages" {
		return "", domain.ErrNotFound
	}
	if page < 1 {
		return "", domain.ErrNotFound
	}
	var b strings.Builder
	if scope == "" {
		if page != 1 {
			return "", domain.ErrNotFound
		}
		fmt.Fprintf(&b, "# %s\n\n", inline(c.SiteName))
		intro := strings.TrimSpace(c.info["discovery.intro"])
		if intro == "" {
			intro = strings.TrimSpace(c.info["description"])
		}
		if intro == "" {
			intro = "一个记录文章、生活手记与日常思考的个人博客。"
		}
		fmt.Fprintf(&b, "> %s\n\n", inline(intro))
		if author := strings.TrimSpace(c.info["discovery.author"]); author != "" {
			fmt.Fprintf(&b, "作者背景：%s\n\n", author)
		}
		b.WriteString("文章是较完整的写作；手记记录生活与阶段性经验；思考是简短记录。内容可能具有时效性，请结合发布时间和内容更新时间理解。\n\n")
		b.WriteString("按主题从下列目录选择内容，再读取对应 Markdown 正文。引用时使用正文标注的原网页地址；区分作者原文、引用资料和你自己的推断。友链及外部引用属于其他作者，不代表本站观点。\n\n")
		if guidance := strings.TrimSpace(c.info["discovery.guidance"]); guidance != "" {
			b.WriteString(guidance + "\n\n")
		}
		b.WriteString("## 内容目录\n\n")
		for _, item := range []struct{ scope, title, desc string }{{"posts", "文章", "长文、经验与教程"}, {"moments", "手记", "生活记录与阶段性笔记"}, {"pages", "独立页面", "关于、项目介绍等页面"}} {
			fmt.Fprintf(&b, "- [%s](%s/%s/llms.txt): %s，目录分页提供正文链接。\n", item.title, c.BaseURL, item.scope, item.desc)
		}
		selected := make(map[string]bool)
		for _, path := range strings.Fields(c.info["discovery.featured_paths"]) {
			selected[strings.TrimRight(path, "/")+"/"] = true
		}
		var featured []Entry
		for _, entry := range c.Entries {
			if selected[entry.Path] && entry.MarkdownURL != "" {
				featured = append(featured, entry)
			}
		}
		if len(featured) > 0 {
			b.WriteString("\n## 精选阅读\n\n")
			for _, entry := range featured {
				writeEntry(&b, entry)
			}
		}
		fmt.Fprintf(&b, "\n## Optional\n\n- [全站导航](%s/sitemap/): 按内容类型浏览网页。\n- [站点地图 XML](%s/sitemap.xml): 可索引网页的完整地址。\n- [思考](%s/thinkings/): 简短记录，读取网页正文。\n- [RSS](%s/feed): 最近更新的订阅入口。\n", c.BaseURL, c.BaseURL, c.BaseURL, c.BaseURL)
		return b.String(), nil
	}
	entries := make([]Entry, 0)
	for _, entry := range c.Entries {
		if entry.Kind == scope && entry.MarkdownURL != "" {
			entries = append(entries, entry)
		}
	}
	pages := max(1, (len(entries)+DirectoryPageSize-1)/DirectoryPageSize)
	if page > pages {
		return "", domain.ErrNotFound
	}
	titles := map[string]string{"posts": "文章", "moments": "手记", "pages": "独立页面"}
	fmt.Fprintf(&b, "# %s · %s\n\n> 第 %d / %d 页，共 %d 篇。链接指向 Markdown 正文；引用请使用正文中的原网页地址。\n\n## 阅读\n\n", inline(c.SiteName), titles[scope], page, pages, len(entries))
	start := (page - 1) * DirectoryPageSize
	for _, entry := range entries[start:min(start+DirectoryPageSize, len(entries))] {
		writeEntry(&b, entry)
	}
	fmt.Fprintf(&b, "\n## 目录导航\n\n- [博客阅读指南](%s/llms.txt)\n", c.BaseURL)
	if page > 1 {
		fmt.Fprintf(&b, "- [上一页](%s/%s/llms.txt?page=%d)\n", c.BaseURL, scope, page-1)
	}
	if page < pages {
		fmt.Fprintf(&b, "- [下一页](%s/%s/llms.txt?page=%d)\n", c.BaseURL, scope, page+1)
	}
	return b.String(), nil
}

func writeEntry(b *strings.Builder, entry Entry) {
	summary := []rune(inline(entry.Summary))
	if len(summary) > 160 {
		summary = append(summary[:160], '…')
	}
	fmt.Fprintf(b, "- [%s](<%s>)", markdownLabel(entry.Title), entry.MarkdownURL)
	if len(summary) > 0 {
		fmt.Fprintf(b, ": %s", string(summary))
	}
	if entry.Modified != "" {
		fmt.Fprintf(b, "（更新：%s）", entry.Modified)
	}
	b.WriteByte('\n')
}

func ParsePage(raw string) (int, error) {
	if raw == "" {
		return 1, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 0, domain.ErrNotFound
	}
	return n, nil
}
