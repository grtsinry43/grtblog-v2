package ws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grtsinry43/grtblog-v2/server/internal/app/article"
	appEvent "github.com/grtsinry43/grtblog-v2/server/internal/app/event"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/content"
	"github.com/grtsinry43/grtblog-v2/server/internal/http/contract"
)

type handlerFunc func(ctx context.Context, event appEvent.Event) error

func (h handlerFunc) Handle(ctx context.Context, event appEvent.Event) error {
	return h(ctx, event)
}

func RegisterArticleUpdateSubscriber(bus appEvent.Bus, manager *Manager) {
	if bus == nil || manager == nil {
		return
	}
	bus.Subscribe(article.ArticleUpdated{}.Name(), handlerFunc(func(ctx context.Context, event appEvent.Event) error {
		updated, ok := event.(article.ArticleUpdated)
		if !ok {
			return nil
		}
		payload := contract.ArticleContentPayload{
			ContentHash: updated.ContentHash,
			Title:       updated.Title,
			LeadIn:      updated.LeadIn,
			TOC:         mapTOCNodes(updated.TOC),
			Content:     updated.Content,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		manager.Broadcast(articleRoomKey(updated.ID), data)
		return nil
	}))
}

func articleRoomKey(id int64) string {
	return fmt.Sprintf("article:%d", id)
}

func mapTOCNodes(nodes []content.TOCNode) []contract.TOCNode {
	result := make([]contract.TOCNode, len(nodes))
	for i, node := range nodes {
		result[i] = contract.TOCNode{
			Name:     node.Name,
			Anchor:   node.Anchor,
			Children: mapTOCNodes(node.Children),
		}
	}
	return result
}
