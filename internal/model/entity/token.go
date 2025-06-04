package entity

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}
