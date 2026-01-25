package clientinfo

import (
	"strings"

	"github.com/ua-parser/uap-go/uaparser"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/comment"
)

type UAParser struct {
	parser *uaparser.Parser
}

func NewUAParser() *UAParser {
	return &UAParser{parser: uaparser.NewFromSaved()}
}

func (p *UAParser) Resolve(userAgent string) comment.ClientInfo {
	ua := strings.TrimSpace(userAgent)
	if ua == "" || p.parser == nil {
		return comment.ClientInfo{}
	}
	client := p.parser.Parse(ua)
	return comment.ClientInfo{
		Platform: strings.TrimSpace(client.Os.Family),
		Browser:  strings.TrimSpace(client.UserAgent.Family),
	}
}

var _ comment.ClientInfoResolver = (*UAParser)(nil)
