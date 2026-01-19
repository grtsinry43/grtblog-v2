package htmlsnapshot

import (
	"context"
	"log"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/article"
	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
)

type handlerFunc func(ctx context.Context, event appEvent.Event) error

func (h handlerFunc) Handle(ctx context.Context, event appEvent.Event) error {
	return h(ctx, event)
}

func RegisterArticleUpdateSubscriber(bus appEvent.Bus, service *Service) {
	if bus == nil || service == nil {
		return
	}
	bus.Subscribe(article.ArticleUpdated{}.Name(), handlerFunc(func(ctx context.Context, event appEvent.Event) error {
		go func() {
			log.Printf("[html-snapshot] trigger from article.updated")
			_ = service.RefreshPostsHTML(context.Background())
		}()
		return nil
	}))
}
