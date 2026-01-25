package article

import (
	"context"
	"time"

	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	appfed "github.com/grtsinry43/grtblog-v2/server/internal/app/federation"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
)

func publishFederationSignals(ctx context.Context, bus appEvent.Bus, article *content.Article, contentBody string) {
	if bus == nil || article == nil {
		return
	}
	mentions, citations := appfed.ParseSignals(contentBody)
	if len(mentions) == 0 && len(citations) == 0 {
		return
	}
	now := time.Now()
	for _, mention := range mentions {
		_ = bus.Publish(ctx, appfed.MentionDetected{
			ArticleID:      article.ID,
			AuthorID:       article.AuthorID,
			Title:          article.Title,
			ShortURL:       article.ShortURL,
			TargetUser:     mention.User,
			TargetInstance: mention.Instance,
			Context:        mention.Context,
			MentionType:    "",
			At:             now,
		})
	}
	for _, citation := range citations {
		_ = bus.Publish(ctx, appfed.CitationDetected{
			ArticleID:      article.ID,
			AuthorID:       article.AuthorID,
			Title:          article.Title,
			ShortURL:       article.ShortURL,
			TargetInstance: citation.Instance,
			TargetPostID:   citation.PostID,
			Context:        citation.Context,
			CitationType:   "",
			At:             now,
		})
	}
}
