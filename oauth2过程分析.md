# oauth2过程分析

### 过程说明

我们可以把oauth2比作「门票兑换系统」

1. 授权码(Code)像兑换劵: 首先，我有一张「兑换券」，但它不能直接进场（即不能直接获得访问权限）

2. 去票务中心换门票：我拿着兑换券去票务中心（授权服务器的/oauth2/token端点），告诉他们这是用来换门票的（授权类型为authorization_code）

3. 票务中心核对信息：票务中心确定你的兑换劵有效，就把一张「门票」（access_token）给你，这张门票是你的进场通行证

4. 凭借门票获得更多信息：有了门票，我可以进入场馆并访问所需要信息（/oauth2/userinfo端点），这张门票为你提供了访问权限

### 简单示例代码

```go
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
	data := "grant_type=authorization_code&code=mock-code" 
	req, _ := http.NewRequest(http.MethodPost, tokenURL, bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") 

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


```

传入access_token参数说明

* grant_type=authorization_code&code=mock-code 表示通过授权码模式换取 access_token。

* Content-Type: application/x-www-form-urlencoded 是指定请求体编码格式，让服务器正确理解传递的参数。

上面代码我发现他使用了httptest去处理mux，这样可以在运行是后台运行一个server，去模拟oauth2的服务器
也知道access token传入POST获取请求时的信息，通过GPT帮我的简化，我已经理解了OAuth的操作流程，感觉自己不仅后面
看他的代码会更容易，而且我感觉我自己可以对接OAuth了
