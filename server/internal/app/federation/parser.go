package federation

import (
	"regexp"
	"strings"
)

type MentionSignal struct {
	User     string
	Instance string
	Context  string
}

type CitationSignal struct {
	Instance string
	PostID   string
	Context  string
}

var (
	mentionPattern  = regexp.MustCompile(`<@([^\s@<>]+)@([^\s<>]+)>`)
	citationPattern = regexp.MustCompile(`<cite:([^|<>]+)\|([^<>]+)>`)
)

// ParseSignals scans article content for mention/citation markers.
func ParseSignals(content string) (mentions []MentionSignal, citations []CitationSignal) {
	if strings.TrimSpace(content) == "" {
		return nil, nil
	}

	mentionSeen := make(map[string]struct{})
	for _, match := range mentionPattern.FindAllStringSubmatchIndex(content, -1) {
		if len(match) < 6 {
			continue
		}
		user := strings.TrimSpace(content[match[2]:match[3]])
		instance := strings.TrimSpace(content[match[4]:match[5]])
		if user == "" || instance == "" {
			continue
		}
		key := user + "@" + instance
		if _, ok := mentionSeen[key]; ok {
			continue
		}
		mentionSeen[key] = struct{}{}
		mentions = append(mentions, MentionSignal{
			User:     user,
			Instance: instance,
			Context:  extractContext(content, match[0], match[1], 80),
		})
	}

	citationSeen := make(map[string]struct{})
	for _, match := range citationPattern.FindAllStringSubmatchIndex(content, -1) {
		if len(match) < 6 {
			continue
		}
		instance := strings.TrimSpace(content[match[2]:match[3]])
		postID := strings.TrimSpace(content[match[4]:match[5]])
		if instance == "" || postID == "" {
			continue
		}
		key := instance + "|" + postID
		if _, ok := citationSeen[key]; ok {
			continue
		}
		citationSeen[key] = struct{}{}
		citations = append(citations, CitationSignal{
			Instance: instance,
			PostID:   postID,
			Context:  extractContext(content, match[0], match[1], 80),
		})
	}

	return mentions, citations
}

func extractContext(content string, start, end, window int) string {
	if window <= 0 {
		return ""
	}
	runes := []rune(content)
	startRune := len([]rune(content[:start]))
	endRune := len([]rune(content[:end]))
	from := startRune - window
	if from < 0 {
		from = 0
	}
	to := endRune + window
	if to > len(runes) {
		to = len(runes)
	}
	return strings.TrimSpace(string(runes[from:to]))
}
