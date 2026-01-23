package webhook

import (
	"fmt"
	"time"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/article"
	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/moment"
	"github.com/grtsinry43/grtblog-v2/server/internal/app/page"
)

var AvailableEventNames = []string{
	article.ArticleCreated{}.Name(),
	article.ArticleUpdated{}.Name(),
	article.ArticlePublished{}.Name(),
	article.ArticleUnpublished{}.Name(),
	article.ArticleDeleted{}.Name(),
	moment.MomentCreated{}.Name(),
	moment.MomentUpdated{}.Name(),
	moment.MomentPublished{}.Name(),
	moment.MomentUnpublished{}.Name(),
	moment.MomentDeleted{}.Name(),
	page.PageCreated{}.Name(),
	page.PageUpdated{}.Name(),
	page.PageDeleted{}.Name(),
}

func IsValidEventName(name string) bool {
	for _, item := range AvailableEventNames {
		if item == name {
			return true
		}
	}
	return false
}

func SampleEvent(name string) (appEvent.Event, error) {
	now := time.Now()
	switch name {
	case article.ArticleCreated{}.Name():
		return article.ArticleCreated{ID: 1, AuthorID: 1, Title: "Sample Article", ShortURL: "sample-article", Published: true, At: now}, nil
	case article.ArticleUpdated{}.Name():
		return article.ArticleUpdated{ID: 1, AuthorID: 1, Title: "Sample Article", ShortURL: "sample-article", Published: true, ContentHash: "hash", LeadIn: nil, TOC: nil, Content: "Sample", At: now}, nil
	case article.ArticlePublished{}.Name():
		return article.ArticlePublished{ID: 1, AuthorID: 1, Title: "Sample Article", ShortURL: "sample-article", At: now}, nil
	case article.ArticleUnpublished{}.Name():
		return article.ArticleUnpublished{ID: 1, AuthorID: 1, Title: "Sample Article", ShortURL: "sample-article", At: now}, nil
	case article.ArticleDeleted{}.Name():
		return article.ArticleDeleted{ID: 1, AuthorID: 1, Title: "Sample Article", ShortURL: "sample-article", At: now}, nil
	case moment.MomentCreated{}.Name():
		return moment.MomentCreated{ID: 1, AuthorID: 1, Title: "Sample Moment", ShortURL: "sample-moment", Published: true, At: now}, nil
	case moment.MomentUpdated{}.Name():
		return moment.MomentUpdated{ID: 1, AuthorID: 1, Title: "Sample Moment", ShortURL: "sample-moment", Published: true, ContentHash: "hash", Summary: "Sample", TOC: nil, Content: "Sample", At: now}, nil
	case moment.MomentPublished{}.Name():
		return moment.MomentPublished{ID: 1, AuthorID: 1, Title: "Sample Moment", ShortURL: "sample-moment", At: now}, nil
	case moment.MomentUnpublished{}.Name():
		return moment.MomentUnpublished{ID: 1, AuthorID: 1, Title: "Sample Moment", ShortURL: "sample-moment", At: now}, nil
	case moment.MomentDeleted{}.Name():
		return moment.MomentDeleted{ID: 1, AuthorID: 1, Title: "Sample Moment", ShortURL: "sample-moment", At: now}, nil
	case page.PageCreated{}.Name():
		return page.PageCreated{ID: 1, Title: "Sample Page", ShortURL: "sample-page", Enabled: true, At: now}, nil
	case page.PageUpdated{}.Name():
		return page.PageUpdated{ID: 1, Title: "Sample Page", ShortURL: "sample-page", Enabled: true, ContentHash: "hash", Description: nil, TOC: nil, Content: "Sample", At: now}, nil
	case page.PageDeleted{}.Name():
		return page.PageDeleted{ID: 1, Title: "Sample Page", ShortURL: "sample-page", At: now}, nil
	default:
		return nil, fmt.Errorf("unknown event: %s", name)
	}
}
