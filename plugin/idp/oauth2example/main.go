package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

// 这是一个简单的Go语言的oauth示例，包括
// 1. Mock 服务器: 模拟 OAuth2 服务端，提供`/oauth2/token` 和`/oauth2/userinfo` 两个端点
// 2. OAuth2 客户端：实现一个简单的客户端，想mock服务器请求`access_token` , 并使用该token 请求用户信息

func main() {
	// 创建一个mock服务器
	mux := http.NewServeMux()

	// 模拟 /oauth2/token 端点
	mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// 响应一个模拟的 access_token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": "mock-access-token",
			"token_type":   "Bearer",
			"expires_in":   "3600",
		})
	})

	// 模拟 /oauth2/userinfo 端点
	mux.HandleFunc("/oauth2/userinfo", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer mock-access-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//响应用户请求
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"sub":   "123456789",
			"name":  "John Doe",
			"email": "john.doe@example.com",
		})
	})

	// 启动mock服务
	server := httptest.NewServer(mux)
	defer server.Close()

	log.Println("Mock server running at:", server.URL)

	// 调用客户端，传mock服务器的地址
	runOAuth2Client(server.URL)
}

func runOAuth2Client(serverURL string) {
	// 1. 请求 access_token
	tokenURL := serverURL + "/oauth2/token"
	data := "grant_type=authorization_code&code=mock-code" // todo 这是什么
	req, _ := http.NewRequest(http.MethodPost, tokenURL, bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // todo 这是什么

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to get token:", err)
		return
	}
	defer resp.Body.Close()

	// 解析 access_token
	var tokenResponse map[string]string
	if err = json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		fmt.Println("Failed to decode token response:", err)
		return
	}

	accessToken := tokenResponse["access_token"]
	fmt.Println("Access Token", accessToken)

	// 2. 使用 access_token 请求用户请求
	userInfoURL := serverURL + "/oauth2/userinfo"
	req, _ = http.NewRequest(http.MethodGet, userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to get user info:", err)
		return
	}
	defer resp.Body.Close()

	// 解析用户信息
	var userInfo map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		fmt.Println("Failed to decode user info response:", err)
		return
	}

	fmt.Println(userInfo)
}
