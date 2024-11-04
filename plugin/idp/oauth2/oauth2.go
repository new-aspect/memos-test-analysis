package oauth2

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
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

// ExchangeToken 这是用没有直接访问权限的授权码(Code)，换取能直接访问的token
// 这里运用了go包里面的oauth2的Exchange
func (p *IdentityProvider) ExchangeToken(ctx context.Context, redirectURL, code string) (string, error) {
	conf := &oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       p.config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   p.config.AuthURL,
			TokenURL:  p.config.TokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	token, err := conf.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange access token %s", err.Error())
	}

	accessToken, ok := token.Extra("access_token").(string)
	if !ok {
		return "", fmt.Errorf("missing \"access_token\" from authorization response")
	}

	return accessToken, nil
}
