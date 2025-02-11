package configs

type Token struct {
	AccessToken string `json:"access_token"`
	AccessUuid  string `json:"access_uuid"`
	AtExpires   int64  `json:"at_expires"`

	RefreshToken     string `json:"refresh_token"`
	RefreshUuid      string `json:"refresh_uuid"`
	RefreshAtExpires int64  `json:"rat_expires"`
}
