package discovery

import (
	"html"
	"regexp"
	"strings"
)

// Attribute syntax mirrors shared/markdown/components.ts, including unquoted
// values and escaped quotes. The export also accepts legacy {key=value} wrappers.
var componentName = regexp.MustCompile(`^(?:\[([^\]\s]+)\]|([^\s\[{]+))`)
var componentAttribute = regexp.MustCompile(`([A-Za-z][\w-]*)\s*=\s*(?:"((?:\\.|[^"\\])*)"|'((?:\\.|[^'\\])*)'|([^\s]+))`)
var componentEscape = regexp.MustCompile(`\\(.)`)
var expandedMention = regexp.MustCompile(`(?s)<fed-mention\s+[^>]*>([^<]*)</fed-mention>`)

func parseExportComponent(info string) (string, map[string]string) {
	info = strings.TrimSpace(info)
	if strings.HasPrefix(info, "component ") {
		info = strings.TrimSpace(strings.TrimPrefix(info, "component "))
	}
	name := componentName.FindStringSubmatch(info)
	attrs := make(map[string]string)
	if name == nil {
		return "", attrs
	}
	raw := strings.TrimSpace(info[len(name[0]):])
	if strings.HasPrefix(raw, "{") && strings.HasSuffix(raw, "}") {
		raw = strings.TrimSpace(raw[1 : len(raw)-1])
	}
	for _, match := range componentAttribute.FindAllStringSubmatchIndex(raw, -1) {
		value := ""
		for group := 2; group <= 4; group++ {
			start, end := match[group*2], match[group*2+1]
			if start < 0 {
				continue
			}
			value = raw[start:end]
			if group != 4 {
				value = componentEscape.ReplaceAllStringFunc(value, func(escaped string) string {
					if escaped == `\n` {
						return "\n"
					}
					return escaped[1:]
				})
			}
			break
		}
		attrs[raw[match[2]:match[3]]] = html.UnescapeString(value)
	}
	return name[1] + name[2], attrs
}

func exportComponent(info string, absolute func(string) string) string {
	name, attrs := parseExportComponent(info)
	title := attrs["title"]
	if title == "" {
		title = attrs["summary"]
	}
	if title == "" {
		title = attrs["caption"]
	}
	link := attrs["href"]
	if link == "" {
		link = attrs["url"]
	}
	if name == "year-card" && title == "" {
		title = "年度总结"
		if attrs["year"] != "" {
			title = attrs["year"] + " 年度总结"
		}
	}
	value := ""
	if link != "" {
		if title == "" {
			title = link
		}
		value = "[" + markdownLabel(title) + "](<" + absolute(link) + ">)"
	} else if title != "" {
		value = "> " + inline(title)
	}
	// Legacy slugs have no defined route in the current YearCard component.
	// Preserve their identity as text instead of manufacturing a broken link.
	if name == "year-card" && link == "" && attrs["slug"] != "" {
		value += "（引用：" + inline(attrs["slug"]) + "）"
	}
	if attrs["desc"] != "" {
		value += "\n\n" + attrs["desc"]
	}
	return value
}
