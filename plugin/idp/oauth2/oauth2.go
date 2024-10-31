package oauth2

import (
	"fmt"
	"new-aspect/memos-test-analysis/store"
)

// IdentityProvider represents an OAuth2 Identity Provider.
type IdentityProvider struct {
	config *store.IdentityProviderOAuth2Config
}

func NewIdentityProvider(config *store.IdentityProviderOAuth2Config) (*IdentityProvider, error) {

	for v, field := range map[string]string{
		config.ClientID:                "clientId",
		config.ClientSecret:            "clientSecret",
		config.TokenURL:                "tokenUrl",
		config.UserInfoURL:             "userInfoUrl",
		config.FieldMapping.Identifier: "fieldMapping.identifier",
	} {
		if v == "" {
			return nil, fmt.Errorf(`the field "%s" is empty but required`, field)
		}
	}

	return &IdentityProvider{
		config: config,
	}, nil
}
