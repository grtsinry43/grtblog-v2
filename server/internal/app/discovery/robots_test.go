package discovery

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRobotsPolicySettings(t *testing.T) {
	settings := testSettings{"public_url": "https://blog.example/"}
	service := NewService(&testRepo{err: errors.New("catalog unavailable")}, settings)
	for _, policy := range []string{"", "   ", "User-agent: *\r\nDisallow: /admin/\r\n\r\nUser-agent: ExampleBot\r\nDisallow: /"} {
		settings["discovery.robots"] = policy
		body, err := service.Robots(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		want := strings.TrimSpace(strings.ReplaceAll(policy, "\r\n", "\n"))
		if want == "" {
			want = "User-agent: *\nDisallow: /admin"
		}
		if body != want+"\n\nSitemap: https://blog.example/sitemap.xml\n" {
			t.Fatalf("unexpected policy: %q", body)
		}
	}
	settings["public_url"] = "https://new.example"
	body, err := service.Robots(context.Background())
	if err != nil || !strings.Contains(body, "Sitemap: https://new.example/sitemap.xml") {
		t.Fatalf("domain update: %q %v", body, err)
	}
}
