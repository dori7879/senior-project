package model

type RefreshTokenForm struct {
	RefreshToken string `json:"refresh_token" form:"max=255"`
}
