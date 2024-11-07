package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"new-aspect/memos-test-analysis/plugin/idp"
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

// 模拟oauth 服务器
// 这个服务器提供两个接口，一个是 /oauth2/token 将授权码Code转换为accessToken
// 另一个是 /oauth2/userinfo 返回用户信息
func newMockServer(t *testing.T, testCoke, testAccessToken string, userInfo []byte) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		// oauth2/token 首先取出来授权码code, 然后判断这个code是否正确
		// 这个是一个POST请求,
		require.Equal(t, http.MethodPost, r.Method)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		vals, err := url.ParseQuery(string(body))
		require.NoError(t, err)

		// 如果code正确，则返回一个accessToken
		var rawIDToken string
		require.Equal(t, testCoke, vals.Get("code"))
		require.Equal(t, "authorization_code", vals.Get("grant_type"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  testAccessToken,
			"token_type":    "Bearer",
			"refresh_token": "test-refresh-token",
			"expires_in":    3600,
			"id_token":      rawIDToken,
		})
		require.NoError(t, err)
	})
	mux.HandleFunc("/oauth2/userinfo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write(userInfo)
		require.NoError(t, err)
	})

	return httptest.NewServer(mux)
}

// 测试整体的逻辑，包含用没有直接访问权限的授权码兑换成有访问权限的token,然后根据token访问用户信息
// 这个过程中需要模拟oauth服务器实现着2次交互，
// 先写发起请求的这2个过程，就是先调用ExchangeToken将授权码兑换成为token
func TestIdentityProvider(t *testing.T) {
	const (
		testClientID    = "test-client-id"
		testCode        = "test-code"
		TestAccessToken = "test-access-token"
		testSubject     = "123456789"
		testName        = "Ask Ning"
		testEmail       = "ask.ning@example.com"
	)

	userInfo, err := json.Marshal(
		map[string]any{
			"sub":   testSubject,
			"name":  testName,
			"email": testEmail,
		},
	)
	require.NoError(t, err)

	s := newMockServer(t, testCode, TestAccessToken, userInfo)

	provider, err := NewIdentityProvider(&store.IdentityProviderOAuth2Config{
		ClientID:     testClientID,
		ClientSecret: testCode,
		TokenURL:     fmt.Sprintf("%s/oauth2/token", s.URL),
		UserInfoURL:  fmt.Sprintf("%s/oauth2/userinfo", s.URL),
		FieldMapping: &store.FieldMapping{
			Identifier:  "sub",
			DisplayName: "name",
			Email:       "email",
		},
	})
	require.NoError(t, err)

	redirectURL := "https://example.com/oauth/callback"
	token, err := provider.ExchangeToken(context.Background(), redirectURL, testCode)
	require.NoError(t, err)

	userInfoResult, err := provider.UserInfo(token)
	require.NoError(t, err)

	// 判断userInfo与我们想要得到的一致
	wantUserInfo := &idp.IdentityProviderUserInfo{
		Identifier:  testSubject,
		DisplayName: testName,
		Email:       testEmail,
	}
	require.Equal(t, userInfoResult, wantUserInfo)
}
