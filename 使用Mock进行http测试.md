# 使用Mockhttp测试
一开始的测试方法

```go
func TestGetHTMLMeta(t *testing.T) {
    tests := []struct {
        urlStr   string
        htmlMeta HTMLMeta
    }{
        {
            urlStr: "https://www.bytebase.com/blog/sql-review-tool-for-devs",
            htmlMeta: HTMLMeta{
                Title:       "The SQL Review Tool for Developers",
                Description: "Reviewing SQL can be somewhat tedious, yet is essential to keep your database fleet reliable. At Bytebase, we are building a developer-first SQL review tool to empower the DevOps system.",
                Image:       "https://www.bytebase.com/static/blog/sql-review-tool-for-devs/dev-fighting-dba.webp",
            },
        },
    }
    for _, test := range tests {
        metadata, err := GetHTMLMeta(test.urlStr)
        require.NoError(t, err)
        require.Equal(t, test.htmlMeta, *metadata)
    }
}
```

有一个问题，就是依赖外部网页测试不稳定，因为网页可能随着时间变化，为了提高测试的稳定性，可以考虑采用Mock HTTP响应，可是使用Go语言的HTTP mock库（例如httptest）来模拟HTTP响应，返回你制定的HTML内容，确保每次测试是可控的

```go
func TestGetHTMLMeta(t *testing.T) {
    // 创建一个模拟的 HTTP服务器
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 设置返回的 HTML 内容
        w.Header().Set("Content-Type", "text/html")
        w.Write([]byte(`
			<!DOCTYPE html>
            <html>
            <head>
                <title>The SQL Review Tool for Developers</title>
                <meta name="description" content="Treat SQL as code...">
                <meta property="og:image" content="https://example.com/image.png">
            </head>
            <body></body>
            </html>
		`))
    }))
    defer server.Close() // 在测试结束时关闭服务器

    // 定义期望的值
    expectedMeta := HTMLMeta{
        Title:       "The SQL Review Tool for Developers",
        Description: "Treat SQL as code...",
        Image:       "https://example.com/image.png",
    }

    // 调用GetHTMLMeta，并验证结果是否符合预期
    metadate, err := GetHTMLMeta(server.URL)
    require.NoError(t, err)
    require.Equal(t, expectedMeta, *metadate)
}
```
在这个测试中，httptest.NewServer 创建了一个模拟的 HTTP 服务器。它会在接收到请求时返回设定的 HTML 响应，从而让 GetHTMLMeta 在访问这个模拟服务器时返回确定的内容，而不依赖外部网页。