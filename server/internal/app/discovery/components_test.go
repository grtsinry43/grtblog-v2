package discovery

import (
	"strings"
	"testing"
)

func TestExportCustomComponents(t *testing.T) {
	for _, test := range []struct{ name, source, want string }{
		{"legacy slug", "::: year-card{slug=2023-summary}\n回顾正文\n:::", "> 年度总结（引用：2023-summary）\n回顾正文\n"},
		{"legacy quoted", `::: year-card{slug="2023-summary" title="我的一年"}`, "> 我的一年（引用：2023-summary）"},
		{"modern", `::: year-card url="/posts/2025/" title="我的 2025"`, "[我的 2025](<https://blog.example/posts/2025/>)"},
		{"unquoted", `::: year-card url=/posts/2025/ title=2025`, "[2025](<https://blog.example/posts/2025/>)"},
		{"escaped", `::: year-card url="/posts/2025/" title="我的 \"2025\""`, "[我的 \"2025\"](<https://blog.example/posts/2025/>)"},
		{"year", `::: year-card year="2024"`, "> 2024 年度总结"},
		{"blockquote", "> ::: year-card{slug=2023}\n> 正文\n> :::", "> > 年度总结（引用：2023）\n> 正文\n> "},
		{"mention", "感谢 <@starnighter@example.com> 的帮助。", "感谢 @starnighter@example.com 的帮助。"},
		{"expanded mention", `<fed-mention user="a" instance="example.com" status="pending">@a@example.com</fed-mention>`, "@a@example.com"},
		{"more unchanged", "正文\n<!--more-->\n后续", "正文\n<!--more-->\n后续"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := CleanMarkdown(test.source, "https://blog.example/posts/current/")
			if got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
}

func TestCustomSyntaxInsideCodeIsUnchanged(t *testing.T) {
	for _, syntax := range []string{`::: year-card{slug=2023}`, `<@starnighter@example.com>`, `<fed-mention user="a" instance="example.com">@a@example.com</fed-mention>`} {
		for _, source := range []string{"`" + syntax + "`", "```markdown\n" + syntax + "\n```", "    " + syntax + "\n"} {
			if got := CleanMarkdown(source, "https://blog.example/"); got != source {
				t.Fatalf("code changed: %q -> %q", source, got)
			}
		}
	}
	source := "::: details summary=\"示例\"\n```markdown\n::: year-card{slug=2023}\n<@starnighter@example.com>\n```\n:::"
	got := CleanMarkdown(source, "https://blog.example/")
	if !strings.Contains(got, "```markdown\n::: year-card{slug=2023}\n<@starnighter@example.com>\n```") {
		t.Fatalf("nested example changed: %s", got)
	}
}
