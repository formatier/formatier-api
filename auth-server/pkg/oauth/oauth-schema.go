package oauth

type OAuthAccountSchema struct {
	Email        string
	AccessToken  string
	RefreshToken string
	ProviderID   string
	Provider     string
}

type AuthCookieSchema struct {
	AccessToken  string `cookie:"access_token"`
	RefreshToken string `cookie:"refresh_token"`
}

type AuthTokenSchema struct {
	Issuer  string   `json:"iss"`
	Subject string   `json:"sub"`
	TokenID string   `json:"jti,omitempty"`
	Email   string   `json:"email"`
	Audince []string `json:"aud"`
	Expiry  uint64   `json:"exp"`
}
