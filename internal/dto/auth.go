package dto

type LoginDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterDto struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Captcha  string `json:"captcha" binding:"required"`
}

// TwitterOAuth2CallbackDto Twitter OAuth2 回调参数
type TwitterOAuth2CallbackDto struct {
	Code  string `json:"code" form:"code" binding:"required"`  // 授权码
	State string `json:"state" form:"state"`                  // 防止CSRF攻击的状态参数
}

// TwitterOAuth2TokenDto Twitter OAuth2 Token响应
type TwitterOAuth2TokenDto struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int    `json:"expires_in"`    // 过期时间(秒)
	TokenType    string `json:"token_type"`    // 令牌类型(Bearer)
}
