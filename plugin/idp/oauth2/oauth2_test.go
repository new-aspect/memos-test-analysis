package oauth2

import (
	"github.com/stretchr/testify/assert"
	"new-aspect/memos-test-analysis/store"
	"testing"
)

func TestNewIdentityProvider(t *testing.T) {
	tests := []struct {
		name        string
		config      *store.IdentityProviderOAuth2Config
		containsErr string
	}{
		{
			name: "no tokenUrl",
			config: &store.IdentityProviderOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				AuthURL:      "",
				TokenURL:     "",
				UserInfoURL:  "https://example.com/api/user",
				FieldMapping: &store.FieldMapping{
					Identifier: "login",
				},
			},
			containsErr: `the field "tokenUrl" is empty but required`,
		},
		{
			name: "no userInfoUrl",
			config: &store.IdentityProviderOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				AuthURL:      "",
				TokenURL:     "https://example.com/token",
				UserInfoURL:  "",
				FieldMapping: &store.FieldMapping{
					Identifier: "login",
				},
			},
			containsErr: `the field "userInfoUrl" is empty but required`,
		},
		{
			name: "no field mapping identifier",
			config: &store.IdentityProviderOAuth2Config{
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				AuthURL:      "",
				TokenURL:     "https://example.com/token",
				UserInfoURL:  "https://example.com/api/user",
				FieldMapping: &store.FieldMapping{
					Identifier: "",
				},
			},
			containsErr: `the field "fieldMapping.identifier" is empty but required`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewIdentityProvider(test.config)
			assert.ErrorContains(t, err, test.containsErr)
		})
	}
}
