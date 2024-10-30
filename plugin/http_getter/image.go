package http_getter

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Image struct {
	Blob      []byte
	Mediatype string
}

func GetImage(urlStr string) (*Image, error) {
	// 验证 URL 是否有效
	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("Invalid URL")
	}

	// 获取图片
	response, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 检查响应的媒体类型
	mediatype, err := getMediatype(response)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(mediatype, "image/") {
		return nil, fmt.Errorf("Expected image, got %s", mediatype)
	}

	// 读取图片数据
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	image := &Image{
		Blob:      bodyBytes,
		Mediatype: mediatype,
	}
	return image, nil
}
