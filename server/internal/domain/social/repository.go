package social

import "context"

type FriendLinkApplicationRepository interface {
	FindByURL(ctx context.Context, url string) (*FriendLinkApplication, error)
	Create(ctx context.Context, app *FriendLinkApplication) error
	Update(ctx context.Context, app *FriendLinkApplication) error
}

type FriendLinkRepository interface {
	FindByURL(ctx context.Context, url string) (*FriendLink, error)
	Create(ctx context.Context, link *FriendLink) error
}
