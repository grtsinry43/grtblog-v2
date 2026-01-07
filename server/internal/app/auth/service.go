package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"

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
	Name       string
	Avatar     string
}

type Service struct {
	users      identity.Repository
	oauthRepo  identity.OAuthProviderRepository
	stateStore StateStore
	manager    *jwt.Manager
	providers  map[string]ExternalProvider
}

func NewService(repo identity.Repository, oauthRepo identity.OAuthProviderRepository, manager *jwt.Manager, stateStore StateStore, authCfg config.AuthConfig) *Service {
	return &Service{
		users:      repo,
		oauthRepo:  oauthRepo,
		stateStore: stateStore,
		manager:    manager,
		providers:  make(map[string]ExternalProvider),
	}
}

func (s *Service) RegisterProvider(provider ExternalProvider) {
	if provider == nil {
		return
	}
	s.providers[strings.ToLower(provider.Name())] = provider
}

type RegisterCmd struct {
	Username string
	Nickname string
	Email    string
	Password string
}

type LoginCmd struct {
	Credential string
	Password   string
}

type LoginResult struct {
	Token  string
	User   identity.User
	Claims *jwt.Claims
}

type UpdateProfileCmd struct {
	UserID   int64
	Nickname string
	Avatar   string
	Email    string
}

type ChangePasswordCmd struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

type AccessInfo struct {
	User identity.User
}

func (s *Service) Register(ctx context.Context, cmd RegisterCmd) (*identity.User, error) {
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
	if total, err := s.users.CountUsers(ctx); err != nil {
		return nil, err
	} else if total == 0 {
		user.IsAdmin = true
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *Service) Login(ctx context.Context, cmd LoginCmd) (*LoginResult, error) {
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
	token, claims, err := s.manager.Generate(user.ID, user.IsAdmin)
	if err != nil {
		return nil, err
	}
	claims.Subject = user.Username
	user.Password = ""
	return &LoginResult{
		Token:  token,
		User:   *user,
		Claims: claims,
	}, nil
}

type OAuthLoginCmd struct {
	Provider string
	Code     string
	State    string
	Redirect string
}

type OAuthAuthorizeResult struct {
	AuthURL       string
	State         string
	CodeChallenge string
}

func (s *Service) LoginWithProvider(ctx context.Context, cmd OAuthLoginCmd) (*LoginResult, error) {
	if s.oauthRepo == nil {
		return nil, ErrProviderNotConfigured
	}
	providerCfg, err := s.oauthRepo.GetByKey(ctx, cmd.Provider)
	if err != nil {
		return nil, err
	}
	if s.stateStore == nil {
		return nil, errors.New("state store not configured")
	}
	stateData, err := s.stateStore.Load(ctx, cmd.State)
	if err != nil {
		return nil, err
	}
	if stateData.Provider != cmd.Provider {
		return nil, errors.New("state/provider mismatch")
	}
	defer s.stateStore.Delete(ctx, cmd.State)

	conf := buildOAuth2Config(providerCfg)
	options := []oauth2.AuthCodeOption{}
	if providerCfg.PKCERequired && stateData.CodeVerifier != "" {
		options = append(options, oauth2.SetAuthURLParam("code_verifier", stateData.CodeVerifier))
	}
	token, err := conf.Exchange(ctx, cmd.Code, options...)
	if err != nil {
		return nil, err
	}

	external, err := fetchExternalIdentity(ctx, providerCfg, token)
	if err != nil {
		return nil, err
	}

	// 映射/注册本地用户
	user, err := s.users.FindByOAuth(ctx, providerCfg.ProviderKey, external.ProviderID)
	if err != nil {
		if errors.Is(err, identity.ErrUserNotFound) {
			user, err = s.registerOAuthUser(ctx, external)
			if err != nil {
				return nil, err
			}
			exp := token.Expiry
			if bindErr := s.users.BindOAuth(ctx, identity.UserOAuth{
				UserID:       user.ID,
				ProviderKey:  providerCfg.ProviderKey,
				OAuthID:      external.ProviderID,
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				ExpiresAt:    &exp,
			}); bindErr != nil {
				return nil, bindErr
			}
		} else {
			return nil, err
		}
	}

	jwtToken, claims, err := s.manager.Generate(user.ID, user.IsAdmin)
	if err != nil {
		return nil, err
	}
	claims.Subject = user.Username
	user.Password = ""
	return &LoginResult{
		Token:  jwtToken,
		User:   *user,
		Claims: claims,
	}, nil
}

// AccessInfo 返回最新的用户、角色与权限信息。
func (s *Service) AccessInfo(ctx context.Context, claims *jwt.Claims) (*AccessInfo, error) {
	if claims == nil {
		return nil, identity.ErrInvalidCredentials
	}
	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return &AccessInfo{
		User: *user,
	}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, cmd UpdateProfileCmd) (*identity.User, error) {
	if cmd.UserID == 0 {
		return nil, identity.ErrInvalidCredentials
	}
	cmd.Nickname = strings.TrimSpace(cmd.Nickname)
	cmd.Email = strings.TrimSpace(cmd.Email)
	cmd.Avatar = strings.TrimSpace(cmd.Avatar)
	updated, err := s.users.UpdateProfile(ctx, cmd.UserID, cmd.Nickname, cmd.Avatar, cmd.Email)
	if err != nil {
		return nil, err
	}
	updated.Password = ""
	return updated, nil
}

// CurrentUser 返回当前用户信息。
func (s *Service) CurrentUser(ctx context.Context, userID int64) (*identity.User, error) {
	if userID == 0 {
		return nil, identity.ErrInvalidCredentials
	}
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *Service) ChangePassword(ctx context.Context, cmd ChangePasswordCmd) error {
	if cmd.UserID == 0 || cmd.NewPassword == "" {
		return identity.ErrInvalidCredentials
	}
	user, err := s.users.FindByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cmd.OldPassword)) != nil {
		return identity.ErrInvalidCredentials
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(cmd.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.users.UpdatePassword(ctx, cmd.UserID, string(hashed))
}

func (s *Service) ListOAuthBindings(ctx context.Context, userID int64) ([]identity.UserOAuthBinding, error) {
	if userID == 0 {
		return nil, identity.ErrInvalidCredentials
	}
	return s.users.ListOAuthBindings(ctx, userID)
}

// IsInitialized 用于判断是否已完成初始化（存在至少一个用户）。
func (s *Service) IsInitialized(ctx context.Context) (bool, error) {
	total, err := s.users.CountUsers(ctx)
	if err != nil {
		return false, err
	}
	return total > 0, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// ListProviders 列出启用的 OAuth 提供方。
func (s *Service) ListProviders(ctx context.Context) ([]identity.OAuthProvider, error) {
	if s.oauthRepo == nil {
		return nil, ErrProviderNotConfigured
	}
	items, err := s.oauthRepo.ListEnabled(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// Authorize 生成授权 URL 和 state（可含 PKCE）。
func (s *Service) Authorize(ctx context.Context, providerKey, redirect string, stateTTL time.Duration) (*OAuthAuthorizeResult, error) {
	if s.oauthRepo == nil || s.stateStore == nil {
		return nil, ErrProviderNotConfigured
	}
	cfg, err := s.oauthRepo.GetByKey(ctx, providerKey)
	if err != nil {
		return nil, err
	}
	oauthCfg := buildOAuth2Config(cfg)

	state, err := GenerateState()
	if err != nil {
		return nil, err
	}
	var codeVerifier, codeChallenge string
	if cfg.PKCERequired {
		codeVerifier, codeChallenge, err = GenerateCodeVerifier()
		if err != nil {
			return nil, err
		}
	}

	if err := s.stateStore.Save(ctx, state, OAuthState{
		Provider:     providerKey,
		Redirect:     redirect,
		CodeVerifier: codeVerifier,
		CreatedAt:    time.Now(),
	}, stateTTL); err != nil {
		return nil, err
	}

	authOpts := []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}
	if cfg.PKCERequired && codeChallenge != "" {
		authOpts = append(authOpts,
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		)
	}
	authURL := oauthCfg.AuthCodeURL(state, authOpts...)
	return &OAuthAuthorizeResult{
		AuthURL:       authURL,
		State:         state,
		CodeChallenge: codeChallenge,
	}, nil
}

func buildOAuth2Config(p *identity.OAuthProvider) *oauth2.Config {
	redirect := p.RedirectURITemplate
	redirect = strings.ReplaceAll(redirect, "{provider}", p.ProviderKey)
	return &oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  p.AuthorizationEndpoint,
			TokenURL: p.TokenEndpoint,
		},
		Scopes:      splitScopes(p.Scopes),
		RedirectURL: redirect,
	}
}

func splitScopes(sc string) []string {
	parts := strings.Fields(sc)
	return parts
}

type ExternalProfile struct {
	ID       string
	Email    string
	Username string
	Name     string
	Avatar   string
}

func fetchExternalIdentity(ctx context.Context, cfg *identity.OAuthProvider, token *oauth2.Token) (*ExternalIdentity, error) {
	profile := ExternalProfile{}
	if cfg.UserinfoEndpoint != "" {
		if err := fetchUserInfo(ctx, cfg.UserinfoEndpoint, token, &profile); err != nil {
			return nil, err
		}
	}
	id := profile.ID
	if id == "" {
		// 回退使用 AccessToken 哈希避免空 ID
		id = token.AccessToken
	}
	return &ExternalIdentity{
		Provider:   cfg.ProviderKey,
		ProviderID: id,
		Email:      profile.Email,
		Username:   firstNonEmpty(profile.Username, profile.Email, profile.Name, id),
		Name:       profile.Name,
		Avatar:     profile.Avatar,
	}, nil
}

func fetchUserInfo(ctx context.Context, endpoint string, token *oauth2.Token, out *ExternalProfile) error {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("userinfo status: %s", resp.Status)
	}
	var raw map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return err
	}
	out.ID = firstString(raw, "sub", "id")
	out.Email = firstString(raw, "email")
	out.Username = firstString(raw, "preferred_username", "username", "login", "name")
	out.Name = firstString(raw, "name")
	out.Avatar = firstString(raw, "avatar_url", "picture")
	return nil
}

func firstString(raw map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := raw[k]; ok {
			if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
				return s
			}
		}
	}
	return ""
}

// registerOAuthUser 根据外部信息注册本地用户。
func (s *Service) registerOAuthUser(ctx context.Context, ext *ExternalIdentity) (*identity.User, error) {
	username := firstNonEmpty(ext.Username, ext.Email, ext.Provider+"_"+ext.ProviderID)
	user := &identity.User{
		Username: username,
		Nickname: ext.Name,
		Email:    ext.Email,
		Avatar:   ext.Avatar,
		IsActive: true,
	}
	if total, err := s.users.CountUsers(ctx); err != nil {
		return nil, err
	} else if total == 0 {
		user.IsAdmin = true
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
