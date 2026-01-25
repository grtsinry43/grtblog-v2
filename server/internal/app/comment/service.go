package comment

import (
	"context"
	"errors"
	"strings"

	domaincomment "github.com/grtsinry43/grtblog-v2/server/internal/domain/comment"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
)

const defaultMaxDepth = 3

type RequestMeta struct {
	IP        string
	UserAgent string
}

type ClientInfo struct {
	Platform string
	Browser  string
}

type ClientInfoResolver interface {
	Resolve(userAgent string) ClientInfo
}

type GeoIPResolver interface {
	Resolve(ip string) string
}

type Service struct {
	repo          domaincomment.CommentRepository
	userRepo      identity.Repository
	clientInfo    ClientInfoResolver
	geoIP         GeoIPResolver
	maxDepthLimit int
}

func NewService(
	repo domaincomment.CommentRepository,
	userRepo identity.Repository,
	clientInfo ClientInfoResolver,
	geoIP GeoIPResolver,
) *Service {
	return &Service{
		repo:          repo,
		userRepo:      userRepo,
		clientInfo:    clientInfo,
		geoIP:         geoIP,
		maxDepthLimit: defaultMaxDepth,
	}
}

type CommentNode struct {
	Comment  *domaincomment.Comment
	Children []*CommentNode
}

func (s *Service) CreateCommentLogin(ctx context.Context, userID int64, cmd CreateCommentLoginCmd, meta RequestMeta) (*domaincomment.Comment, error) {
	if err := s.ensureContentValid(cmd.Content); err != nil {
		return nil, err
	}
	if err := s.ensureAreaExists(ctx, cmd.AreaID); err != nil {
		return nil, err
	}

	if cmd.ParentID != nil {
		if err := s.ensureParentValid(ctx, cmd.AreaID, *cmd.ParentID); err != nil {
			return nil, err
		}
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	nickname := strings.TrimSpace(user.Nickname)
	if nickname == "" {
		nickname = strings.TrimSpace(user.Username)
	}
	nicknamePtr := toPtr(nickname)
	emailPtr := toPtr(strings.TrimSpace(user.Email))

	commentEntity := &domaincomment.Comment{
		AreaID:   cmd.AreaID,
		Content:  strings.TrimSpace(cmd.Content),
		AuthorID: &user.ID,
		NickName: nicknamePtr,
		Email:    emailPtr,
		Website:  nil,
		IsOwner:  user.IsAdmin,
		IsAuthor: user.IsAdmin,
		IsFriend: false,
		IsViewed: user.IsAdmin,
		IsTop:    false,
		ParentID: cmd.ParentID,
	}
	s.applyRequestMeta(commentEntity, meta)

	if err := s.repo.Create(ctx, commentEntity); err != nil {
		return nil, err
	}
	return commentEntity, nil
}

func (s *Service) CreateCommentVisitor(ctx context.Context, cmd CreateCommentVisitorCmd, meta RequestMeta) (*domaincomment.Comment, error) {
	if err := s.ensureContentValid(cmd.Content); err != nil {
		return nil, err
	}
	if err := s.ensureAreaExists(ctx, cmd.AreaID); err != nil {
		return nil, err
	}
	if cmd.ParentID != nil {
		if err := s.ensureParentValid(ctx, cmd.AreaID, *cmd.ParentID); err != nil {
			return nil, err
		}
	}

	nickname := strings.TrimSpace(cmd.NickName)
	email := strings.TrimSpace(cmd.Email)
	website := strings.TrimSpace(toValue(cmd.Website))

	commentEntity := &domaincomment.Comment{
		AreaID:   cmd.AreaID,
		Content:  strings.TrimSpace(cmd.Content),
		AuthorID: nil,
		NickName: toPtr(nickname),
		Email:    toPtr(email),
		Website:  toPtr(website),
		IsOwner:  false,
		IsAuthor: false,
		IsFriend: false,
		IsViewed: false,
		IsTop:    false,
		ParentID: cmd.ParentID,
	}
	s.applyRequestMeta(commentEntity, meta)

	if err := s.repo.Create(ctx, commentEntity); err != nil {
		return nil, err
	}
	return commentEntity, nil
}

func (s *Service) ListCommentTree(ctx context.Context, areaID int64) ([]*CommentNode, error) {
	if err := s.ensureAreaExists(ctx, areaID); err != nil {
		return nil, err
	}
	items, err := s.repo.ListByAreaID(ctx, areaID)
	if err != nil {
		return nil, err
	}
	return buildCommentTree(items), nil
}

func (s *Service) applyRequestMeta(commentEntity *domaincomment.Comment, meta RequestMeta) {
	ip := strings.TrimSpace(meta.IP)
	if ip != "" {
		commentEntity.IP = &ip
	}
	if s.clientInfo != nil {
		info := s.clientInfo.Resolve(meta.UserAgent)
		if strings.TrimSpace(info.Platform) != "" {
			commentEntity.Platform = toPtr(info.Platform)
		}
		if strings.TrimSpace(info.Browser) != "" {
			commentEntity.Browser = toPtr(info.Browser)
		}
	}
	if s.geoIP != nil {
		location := strings.TrimSpace(s.geoIP.Resolve(ip))
		if location != "" {
			commentEntity.Location = &location
		}
	}
}

func (s *Service) ensureContentValid(content string) error {
	if strings.TrimSpace(content) == "" {
		return domaincomment.ErrCommentContentEmpty
	}
	return nil
}

func (s *Service) ensureAreaExists(ctx context.Context, areaID int64) error {
	if areaID <= 0 {
		return domaincomment.ErrCommentAreaNotFound
	}
	area, err := s.repo.GetAreaByID(ctx, areaID)
	if err != nil {
		if errors.Is(err, domaincomment.ErrCommentAreaNotFound) {
			return err
		}
		return err
	}
	if area.IsClosed {
		return domaincomment.ErrCommentAreaClosed
	}
	return nil
}

func (s *Service) ensureParentValid(ctx context.Context, areaID int64, parentID int64) error {
	parent, err := s.repo.FindByID(ctx, parentID)
	if err != nil {
		if errors.Is(err, domaincomment.ErrCommentNotFound) {
			return domaincomment.ErrCommentParentNotFound
		}
		return err
	}
	if parent.AreaID != areaID {
		return domaincomment.ErrCommentParentNotFound
	}

	chainLength := 1
	current := parent
	for current.ParentID != nil {
		if chainLength+1 >= s.maxDepthLimit {
			return domaincomment.ErrCommentTooDeep
		}
		next, err := s.repo.FindByID(ctx, *current.ParentID)
		if err != nil {
			if errors.Is(err, domaincomment.ErrCommentNotFound) {
				return domaincomment.ErrCommentParentNotFound
			}
			return err
		}
		chainLength++
		current = next
	}
	if chainLength+1 > s.maxDepthLimit {
		return domaincomment.ErrCommentTooDeep
	}
	return nil
}

func buildCommentTree(items []*domaincomment.Comment) []*CommentNode {
	nodes := make(map[int64]*CommentNode, len(items))
	for _, item := range items {
		nodes[item.ID] = &CommentNode{Comment: item}
	}

	var roots []*CommentNode
	for _, item := range items {
		node := nodes[item.ID]
		if item.ParentID != nil {
			if parent, ok := nodes[*item.ParentID]; ok {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	return roots
}

func toPtr(val string) *string {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func toValue(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}
