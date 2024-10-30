package http_getter

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetImage(t *testing.T) {
	// 创建模拟服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置 Content-Type 头部，模拟图片
		w.Header().Set("Content-Type", "image/webp")
		w.Write([]byte{0x1, 0x2, 0x3, 0x4}) // 写入模拟的图片数据
	}))
	defer server.Close() // 在测试结束后关闭服务器

	tests := []struct {
		urlStr       string
		expectedType string
	}{
		{
			urlStr:       server.URL,
			expectedType: "image/webp",
		},
	}

	for _, test := range tests {
		image, err := GetImage(test.urlStr)
		require.NoError(t, err)
		require.Equal(t, test.expectedType, image.Mediatype)
		require.NotEmpty(t, image.Blob)
	}
}
