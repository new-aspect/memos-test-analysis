package http_getter

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
