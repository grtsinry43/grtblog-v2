package article

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

type heading struct {
	level int
	text  string
}

func generateTOC(markdown string) []content.TOCNode {
	headings := extractHeadings(markdown)
	if len(headings) == 0 {
		return nil
	}

	minLevel := headings[0].level
	for _, h := range headings[1:] {
		if h.level < minLevel {
			minLevel = h.level
		}
	}

	var roots []*content.TOCNode
	levelMap := make(map[int][]*content.TOCNode)
	anchorIndex := 1
	for _, h := range headings {
		node := &content.TOCNode{
			Name:   h.text,
			Anchor: "article-md-title-" + strconv.Itoa(anchorIndex),
		}
		anchorIndex++

		if h.level == minLevel {
			roots = append(roots, node)
		} else {
			parentList := levelMap[h.level-1]
			if len(parentList) > 0 {
				parent := parentList[len(parentList)-1]
				parent.Children = append(parent.Children, content.TOCNode{
					Name:   node.Name,
					Anchor: node.Anchor,
				})
				node = &parent.Children[len(parent.Children)-1]
			} else {
				roots = append(roots, node)
			}
		}
		levelMap[h.level] = append(levelMap[h.level], node)
	}

	result := make([]content.TOCNode, len(roots))
	for i, node := range roots {
		result[i] = *node
	}
	return result
}

func extractHeadings(markdown string) []heading {
	scanner := bufio.NewScanner(strings.NewReader(markdown))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	inCodeBlock := false
	var headings []heading
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		level := 0
		for level < len(line) && line[level] == '#' {
			level++
		}
		if level == 0 || level > 6 {
			continue
		}
		text := strings.TrimSpace(line[level:])
		if text == "" {
			continue
		}
		headings = append(headings, heading{
			level: level,
			text:  text,
		})
	}
	return headings
}

func buildSummary(summary, content string) string {
	if strings.TrimSpace(summary) != "" {
		return summary
	}
	return truncateRunes(content, 200)
}

func truncateRunes(input string, limit int) string {
	runes := []rune(input)
	if len(runes) <= limit {
		return input
	}
	return string(runes[:limit])
}
