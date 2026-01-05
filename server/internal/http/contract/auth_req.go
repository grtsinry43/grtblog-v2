package contract

// RegisterReq 注册请求。
type RegisterReq struct {
	Username       string `json:"username"`
	Nickname       string `json:"nickname"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	TurnstileToken string `json:"turnstileToken"`
}

// LoginReq 登录请求。
type LoginReq struct {
	Credential     string `json:"credential"` // username or email
	Password       string `json:"password"`
	TurnstileToken string `json:"turnstileToken"`
}

// UpdateProfileReq 更新资料请求。
type UpdateProfileReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

// ChangePasswordReq 修改密码请求。
type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
