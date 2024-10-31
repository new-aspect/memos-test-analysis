package store

type IdentityProviderOAuth2Config struct {
	ClientID     string        `json:"clientId"`
	ClientSecret string        `json:"clientSecret"`
	AuthURL      string        `json:"authUrl"`
	TokenURL     string        `json:"tokenUrl"`
	UserInfoURL  string        `json:"userInfoUrl"`
	Scopes       []string      `json:"scopes"`
	FieldMapping *FieldMapping `json:"fieldMapping"`
}

type FieldMapping struct {
	Identifier  string `json:"identifier"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}
