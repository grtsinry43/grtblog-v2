package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/grtsinry43/grtblog-v2/server/internal/config"
	"github.com/grtsinry43/grtblog-v2/server/internal/domain/identity"
	"github.com/grtsinry43/grtblog-v2/server/internal/security/jwt"
)

var (
	ErrProviderNotConfigured = errors.New("oauth provider not configured")
)

// ExternalProvider 用于未来扩展 OAuth/OIDC, 当前仅定义接口。
type ExternalProvider interface {
	Name() string
	Exchange(ctx context.Context, code string) (*ExternalIdentity, error)
}

type ExternalIdentity struct {
	Provider   string
	ProviderID string
	Email      string
	Username   string
	Avatar     string
}

type Service struct {
	users        identity.Repository
	manager      *jwt.Manager
	defaultRoles []string
	providers    map[string]ExternalProvider
}

func NewService(repo identity.Repository, manager *jwt.Manager, authCfg config.AuthConfig) *Service {
	return &Service{
		users:        repo,
		manager:      manager,
		defaultRoles: authCfg.DefaultRoles,
		providers:    make(map[string]ExternalProvider),
	}
}

func (s *Service) RegisterProvider(provider ExternalProvider) {
	if provider == nil {
		return
	}
	s.providers[strings.ToLower(provider.Name())] = provider
}

type RegisterCommand struct {
	Username string
	Nickname string
	Email    string
	Password string
}

type LoginCommand struct {
	Credential string
	Password   string
}

type LoginResult struct {
	Token  string
	User   identity.User
	Claims *jwt.Claims
	Roles  []string
}

func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (*identity.User, error) {
	cmd.Username = strings.TrimSpace(cmd.Username)
	if cmd.Username == "" || cmd.Password == "" {
		return nil, identity.ErrInvalidCredentials
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &identity.User{
		Username: cmd.Username,
		Nickname: firstNonEmpty(cmd.Nickname, cmd.Username),
		Email:    strings.TrimSpace(cmd.Email),
		Password: string(hashed),
		IsActive: true,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	if err := s.users.AssignRoles(ctx, user.ID, s.defaultRoles); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *Service) Login(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	if cmd.Credential == "" || cmd.Password == "" {
		return nil, identity.ErrInvalidCredentials
	}
	user, err := s.users.FindByCredential(ctx, cmd.Credential)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cmd.Password)) != nil {
		return nil, identity.ErrInvalidCredentials
	}
	roles, err := s.users.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	perms, err := s.users.GetPermissions(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	token, claims, err := s.manager.Generate(user.ID, roles, perms)
	if err != nil {
		return nil, err
	}
	claims.Subject = user.Username
	user.Password = ""
	return &LoginResult{
		Token:  token,
		User:   *user,
		Claims: claims,
		Roles:  roles,
	}, nil
}

type OAuthLoginCommand struct {
	Provider string
	Code     string
}

func (s *Service) LoginWithProvider(ctx context.Context, cmd OAuthLoginCommand) (*LoginResult, error) {
	provider, ok := s.providers[strings.ToLower(cmd.Provider)]
	if !ok {
		return nil, ErrProviderNotConfigured
	}
	external, err := provider.Exchange(ctx, cmd.Code)
	if err != nil {
		return nil, err
	}
	// TODO: 将 external identity 映射 / 注册到本地用户体系
	_ = external
	return nil, errors.New("oauth login not implemented yet")
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
