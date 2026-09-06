package handler

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	app "github.com/grtsinry43/grtblog-v2/server/internal/app/discovery"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/sysconfig"
	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence"
	"github.com/grtsinry43/grtblog-v2/server/internal/infra/persistence/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Runs against an explicitly supplied disposable PostgreSQL database. Each run
// owns an isolated schema and drops only that schema when finished.
func TestDiscoveryPostgresLifecycle(t *testing.T) {
	dsn := os.Getenv("DISCOVERY_TEST_DSN")
	if dsn == "" {
		t.Skip("set DISCOVERY_TEST_DSN to a disposable PostgreSQL database")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })
	schema := fmt.Sprintf("discovery_test_%d", time.Now().UnixNano())
	if err := db.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Exec("SET search_path TO public"); db.Exec("DROP SCHEMA " + schema + " CASCADE") })
	if err := db.Exec("SET search_path TO " + schema).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Article{}, &model.Moment{}, &model.Page{}, &model.ArticleCategory{}, &model.MomentColumn{}, &model.Album{}, &model.Photo{}, &model.User{}, &model.SysConfig{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("CREATE UNIQUE INDEX ON sys_config(config_key)").Error; err != nil {
		t.Fatal(err)
	}
	migration, err := os.ReadFile("../../../migrations/0070_add_discovery_settings.sql")
	if err != nil {
		t.Fatal(err)
	}
	up := strings.Split(string(migration), "-- +goose Down")[0]
	if err := db.Exec(up).Error; err != nil {
		t.Fatal(err)
	}
	// Verify the migration is safe to run twice and produces discoverable UI metadata.
	if err := db.Exec(up).Error; err != nil {
		t.Fatal(err)
	}
	var count int64
	db.Model(&model.SysConfig{}).Where("group_path = ?", "site/discovery").Count(&count)
	if count != 5 {
		t.Fatalf("settings count %d", count)
	}
	for key, value := range map[string]string{"site.public_url": "https://blog.example", "site.website_name": "测试博客"} {
		if err := db.Create(&model.SysConfig{ConfigKey: key, Value: value}).Error; err != nil {
			t.Fatal(err)
		}
	}
	created := time.Date(2026, 9, 5, 18, 0, 0, 0, time.UTC)
	author := model.User{Username: "writer", Nickname: "测试作者", IsActive: true}
	if err := db.Create(&author).Error; err != nil {
		t.Fatal(err)
	}
	article := model.Article{Title: "公开文章 & XML", ShortURL: "public", Summary: "摘要", Content: "正文 ![图](/uploads/test.png)", TOC: []byte("[]"), AuthorID: author.ID, IsPublished: true, CreatedAt: created, ContentUpdatedAt: created}
	draft := article
	draft.ID = 0
	draft.ShortURL = "draft"
	draft.IsPublished = false
	draft.Content = "PRIVATE_DRAFT"
	for _, item := range []*model.Article{&article, &draft} {
		if err := db.Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
	moment := model.Moment{Title: "手记", ShortURL: "note", Content: "手记正文", TOC: []byte("[]"), AuthorID: author.ID, IsPublished: true, CreatedAt: created, ContentUpdatedAt: created}
	if err := db.Create(&moment).Error; err != nil {
		t.Fatal(err)
	}
	for _, item := range []model.Page{{Title: "关于", ShortURL: "about", Content: "关于正文", IsEnabled: true, TOC: []byte("[]"), CreatedAt: created, ContentUpdatedAt: created}, {Title: "禁用", ShortURL: "disabled", Content: "PRIVATE_PAGE", TOC: []byte("[]"), CreatedAt: created, ContentUpdatedAt: created}, {Title: "内置", ShortURL: "posts", Content: "BUILTIN_BODY", IsEnabled: true, IsBuiltin: true, TOC: []byte("[]"), CreatedAt: created, ContentUpdatedAt: created}} {
		if err := db.Create(&item).Error; err != nil {
			t.Fatal(err)
		}
	}
	album := model.Album{Title: "相册", ShortURL: "photos", IsPublished: true, AuthorID: author.ID}
	if err := db.Create(&album).Error; err != nil {
		t.Fatal(err)
	}
	photo := model.Photo{AlbumID: &album.ID, URL: "/uploads/photo.jpg"}
	if err := db.Create(&photo).Error; err != nil {
		t.Fatal(err)
	}
	sys := sysconfig.NewService(persistence.NewSysConfigRepository(db), config.TurnstileConfig{}, nil)
	if err := sys.UpdateWebsiteInfoByKey(context.Background(), "discovery.intro", "来自后台的博客介绍"); err != nil {
		t.Fatal(err)
	}
	if err := sys.UpdateWebsiteInfoByKey(context.Background(), "discovery.featured_paths", "/posts/public/\n/posts/draft/"); err != nil {
		t.Fatal(err)
	}
	h := NewDiscoveryHandler(app.NewService(persistence.NewDiscoveryRepository(db), sys))
	web := fiber.New()
	web.Get("/catalog", h.Catalog)
	web.Get("/resource/*", h.Resource)
	request := func(method, path, etag string) (int, string, string, string) {
		t.Helper()
		req := httptest.NewRequest(method, path, nil)
		if etag != "" {
			req.Header.Set("If-None-Match", etag)
		}
		res, err := web.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		return res.StatusCode, string(body), res.Header.Get("ETag"), res.Header.Get("Content-Type")
	}
	status, xml, etag, contentType := request("GET", "/resource/sitemap.xml", "")
	if status != 200 || !strings.HasPrefix(contentType, "application/xml") || !strings.Contains(xml, "/posts/public/") || !strings.Contains(xml, fmt.Sprintf("/albums/photos/photo/%d/", photo.ID)) || strings.Contains(xml, "draft") || strings.Contains(xml, "disabled") {
		t.Fatalf("%d %s %s", status, contentType, xml)
	}
	if code, _, _, _ := request("GET", "/resource/sitemap.xml", "W/"+etag); code != 304 {
		t.Fatalf("conditional GET: %d", code)
	}
	if code, body, _, _ := request("HEAD", "/resource/sitemap.xml", ""); code != 200 || body != "" {
		t.Fatalf("HEAD: %d %s", code, body)
	}
	if code, body, _, _ := request("GET", "/resource/llms.txt", ""); code != 200 || !strings.Contains(body, "来自后台的博客介绍") || !strings.Contains(body, "/posts/public/index.md") || strings.Contains(body, "draft") {
		t.Fatalf("LLMS %d %s", code, body)
	}
	if code, body, _, _ := request("GET", "/resource/posts/public/index.md", ""); code != 200 || !strings.Contains(body, "测试作者") || !strings.Contains(body, "正文") || !strings.Contains(body, "https://blog.example/uploads/test.png") {
		t.Fatalf("Markdown %d %s", code, body)
	}
	for _, path := range []string{"posts/draft/index.md", "disabled/index.md", "posts/index.md", "moments/2026/09/05/note/index.md", "posts/llms.txt?page=99", "sitemaps/0.xml"} {
		if code, body, _, _ := request("GET", "/resource/"+path, ""); code != 404 || strings.Contains(body, "PRIVATE") {
			t.Fatalf("%s: %d %s", path, code, body)
		}
	}
	if code, body, _, _ := request("GET", "/resource/moments/2026/09/06/note/index.md", ""); code != 200 || !strings.Contains(body, "手记正文") {
		t.Fatalf("moment: %d %s", code, body)
	}
	if err := db.Model(&article).Update("is_published", false).Error; err != nil {
		t.Fatal(err)
	}
	if code, body, _, _ := request("GET", "/resource/sitemap.xml", etag); code != 200 || strings.Contains(body, "/posts/public/") {
		t.Fatalf("stale sitemap %d %s", code, body)
	}
	if code, _, _, _ := request("GET", "/resource/posts/public/index.md", ""); code != 404 {
		t.Fatalf("withdrawn document %d", code)
	}
	if _, body, _, _ := request("GET", "/resource/llms.txt", ""); strings.Contains(body, "/posts/public/") {
		t.Fatal("withdrawn featured entry")
	}
	if err := db.Model(&article).Update("is_published", true).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Delete(&article).Error; err != nil {
		t.Fatal(err)
	}
	if code, _, _, _ := request("GET", "/resource/posts/public/index.md", ""); code != 404 {
		t.Fatalf("deleted document %d", code)
	}
	if err := db.Model(&album).Update("is_published", false).Error; err != nil {
		t.Fatal(err)
	}
	if _, body, _, _ := request("GET", "/resource/sitemap.xml", ""); strings.Contains(body, "/albums/photos/") {
		t.Fatal("private album/photos indexed")
	}
	// A query failure must not masquerade as an empty successful sitemap.
	if err := db.Exec("ALTER TABLE article RENAME TO unavailable_article").Error; err != nil {
		t.Fatal(err)
	}
	if code, _, _, _ := request("GET", "/resource/sitemap.xml", ""); code != 503 {
		t.Fatalf("outage status %d", code)
	}
}
