package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"new-aspect/memos-test-analysis/plugin/idp"
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

// UserInfo return the parsed user information using the given OAuth2 token.
// 也就是这个方法查询userinfo的http请求，然后将返回的请求结构整理成为IdentifyUserInfo结构
//
// 实现方式，使用http发出一个请求头包含token的userInfoURL的请求
// 将返回结果json序列化以后根据p.config.FieldMapping.Identifier等3个字段拼接出来IdentifyUserInfo结构
//
// 疑问，就是我好奇他为什么返回的时候选择idp包下面的一个结构，而不是直接返回一个在oauth2包下面的结构
// 简单来说是为了解耦，为了清晰的划分职责，为了方便后面拓展不同类似的身份认证的时候减少耦合
func (p *IdentityProvider) UserInfo(token string) (*idp.IdentityProviderUserInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, p.config.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to new http request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user information")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body %s", err)
	}

	var claims map[string]any
	err = json.Unmarshal(body, &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body %s", err)
	}

	userInfo := &idp.IdentityProviderUserInfo{}
	if v, ok := claims[p.config.FieldMapping.Identifier].(string); ok {
		userInfo.Identifier = v
	}
	if userInfo.Identifier == "" {
		return nil, fmt.Errorf("the field %q is not found in claims or has empty value", p.config.FieldMapping.Identifier)
	}

	// Best effort to map optional fields
	if p.config.FieldMapping.DisplayName != "" {
		if v, ok := claims[p.config.FieldMapping.DisplayName].(string); ok {
			userInfo.DisplayName = v
		}
	}
	if userInfo.DisplayName == "" {
		userInfo.DisplayName = userInfo.Identifier
	}
	if p.config.FieldMapping.Email != "" {
		if v, ok := claims[p.config.FieldMapping.Email].(string); ok {
			userInfo.Email = v
		}
	}
	return userInfo, nil
}
