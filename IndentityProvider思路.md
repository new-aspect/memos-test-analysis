# NewIdentityProvider 思路

```go
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
```

简单来说这个 NewIdentityProvider 在创建之前做了一个判断，判断是否所有字符串
都有值，他没有简单的用 if 判断每一个变量，因为这样子写下来就有太多重复代码了

```go
func NewIdentityProvider(config *store.IdentityProviderOAuth2Config) (*IdentityProvider, error) {
	if config.ClientID == "" {
		return nil, fmt.Errorf(`the field "%s" is empty but required`, "clientId")
	}
	if config.ClientSecret == "" {
		return nil, fmt.Errorf(`the field "%s" is empty but required`, "clientSecret")
	}
	if config.TokenURL == "" {
		return nil, fmt.Errorf(`the field "%s" is empty but required`, "tokenUrl")
	}
	if config.UserInfoURL == "" {
		return nil, fmt.Errorf(`the field "%s" is empty but required`, "userInfoUrl")
	}
	if config.FieldMapping.Identifier == "" {
		return nil, fmt.Errorf(`the field "%s" is empty but required`, "fieldMapping.identifier")
	}

	return &IdentityProvider{
		config: config,
	}, nil
}
```

所以通过一个map和遍历避免了上面的这种重复代码

看测试案例，它首先创建了一个oauth2，我们的授权码(Code), 这个授权码是没有访问权限的，然后用这个oauth2根据授权码(Code)换取有访问权限的Token,
拿到了这个访问权限以后，我们尝试拿到我们想要的信息(/oauth2/userinfo端点)，

然后我就突然有感想我的测试文件就是在描述这个小的单元是怎么工作的。然后在这个过程中我有一个疑问，就是为什么作者选择集成测试，而不是将ExchangeToken和UserInfo
分开做单独的单元测试。我看到GPT分析打动我的一点是这么做集成测试可以减少测试对实现细节的依赖，作者认为Oauth2流程中的ExchangeToken和UserInfo本身是耦合度高的逻辑，
将其分开测试可能使得测试对实现细节的依赖性更强。这样当逻辑变动时，不需要过多的调整单元测试，减少了维护工作量。

