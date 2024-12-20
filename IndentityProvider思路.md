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



### 为什么UserInfo返回的时候选择idp包下面的一个结构，而不是直接返回一个在oauth2包下面的结构

这是一个设计选择，主要是为了让Oauth2包与idp保持解耦，原因如下

1. 模块化和职责分离: 这种设计使得每个包只关心它自己的职责，避免了不同模块之间的紧密耦合。

2. 统一的数据格式: OAuth2、LDAP、SAML等用户信息的格式和字段可能有所不同,通过使用idp.IdentityProviderUserInfo作为统一的返回类型，系统能够方便地扩展其他身份验证方式，而不需要修改oauth2包的内部逻辑。比如，如果后续添加了其他身份验证方式（比如LDAP），可以直接使用idp.IdentityProviderUserInfo结构返回用户信息。

3. 易于扩展和维护: 如果将用户信息直接返回为oauth2包中的结构体，可能会导致oauth2包承担更多职责，未来扩展时也更难以管理。通过将用户信息单独放在idp包中，代码的扩展性和维护性都得到了增强。如果以后需要在IdentityProviderUserInfo中添加更多字段（如phone number、address等），只需要修改idp包，而不需要改动oauth2包中的代码。

4. 统一的用户信息模型：有时候，一个OAuth2身份提供者可能会返回的用户信息与该应用的用户信息模型不完全一致。通过idp.IdentityProviderUserInfo，可以在oauth2包内部进行转换处理后，将其统一成应用需要的格式，而不暴露给外部使用者。这样做确保了oauth2包的返回数据符合应用的需求，而不需要外部调用者关心内部的转换细节。
