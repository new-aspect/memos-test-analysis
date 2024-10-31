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