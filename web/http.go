// Package web 网络处理相关
package web

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"

	trshttp "github.com/fumiama/terasu/http"
)

// NewDefaultClient ...
func NewDefaultClient() *http.Client {
	cp := trshttp.DefaultClient
	return &cp
}

// NewTLS12Client ...
func NewTLS12Client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialTLSContext: trshttp.DefaultClient.Transport.(*http.Transport).DialTLSContext,
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS12,
			},
		},
	}
}

// NewPixivClient P站特殊客户端
func NewPixivClient() *http.Client {
	return NewTLS12Client()
}

// RequestDataWith 使用自定义请求头获取数据
func RequestDataWith(client *http.Client, url, method, referer, ua string, body io.Reader) (data []byte, err error) {
	// 提交请求
	var request *http.Request
	request, err = http.NewRequest(method, url, body)
	if err == nil {
		// 增加header选项
		if referer != "" {
			request.Header.Add("Referer", referer)
		}
		if ua != "" {
			request.Header.Add("User-Agent", ua)
		}
		var response *http.Response
		response, err = client.Do(request)
		if err == nil {
			if response.StatusCode != http.StatusOK {
				s := fmt.Sprintf("status code: %d", response.StatusCode)
				err = errors.New(s)
				return
			}
			data, err = io.ReadAll(response.Body)
			response.Body.Close()
		}
	}
	return
}

// RequestDataWithHeaders 使用自定义请求头获取数据
func RequestDataWithHeaders(client *http.Client, url, method string, setheaders func(*http.Request) error, body io.Reader) (data []byte, err error) {
	// 提交请求
	var request *http.Request
	request, err = http.NewRequest(method, url, body)
	if err == nil {
		// 增加header选项
		err = setheaders(request)
		if err != nil {
			return
		}
		var response *http.Response
		response, err = client.Do(request)
		if err != nil {
			return
		}
		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusPartialContent {
			s := fmt.Sprintf("status code: %d", response.StatusCode)
			err = errors.New(s)
			return
		}
		data, err = io.ReadAll(response.Body)
		response.Body.Close()
	}
	return
}

// GetData 获取数据
func GetData(url string) (data []byte, err error) {
	var response *http.Response
	response, err = trshttp.Get(url)
	if err != nil {
		response, err = http.Get(url)
	}
	if err == nil {
		if response.StatusCode != http.StatusOK {
			s := fmt.Sprintf("status code: %d", response.StatusCode)
			err = errors.New(s)
			return
		}
		data, err = io.ReadAll(response.Body)
		response.Body.Close()
	}
	return
}

// PostData 获取数据
func PostData(url, contentType string, body io.Reader) (data []byte, err error) {
	var response *http.Response
	response, err = trshttp.Post(url, contentType, body)
	if err != nil {
		response, err = http.Post(url, contentType, body)
	}
	if err == nil {
		if response.StatusCode != http.StatusOK {
			s := fmt.Sprintf("status code: %d", response.StatusCode)
			err = errors.New(s)
			return
		}
		data, err = io.ReadAll(response.Body)
		response.Body.Close()
	}
	return
}

// HeadRequestURL 获取跳转后的链接
func HeadRequestURL(u string) (newu string, err error) {
	var data *http.Response
	data, err = trshttp.Head(u)
	if err != nil {
		data, err = http.Head(u)
	}
	if err != nil {
		return "", err
	}
	_ = data.Body.Close()
	return data.Request.URL.String(), nil
}
