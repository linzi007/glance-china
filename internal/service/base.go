package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// BaseClient 基础客户端实现
type BaseClient struct {
	name      string
	baseURL   string
	timeout   time.Duration
	headers   map[string]string
	client    *http.Client
}

func (b *BaseClient) GetName() string {
	return b.name
}

func (b *BaseClient) GetBaseURL() string {
	return b.baseURL
}

func (b *BaseClient) IsHealthy(ctx context.Context) bool {
	req := &APIRequest{
		Method:  "GET",
		Path:    "/health",
		Timeout: 5 * time.Second,
	}
	
	_, err := b.Request(ctx, req)
	return err == nil
}

func (b *BaseClient) SetRateLimit(requests int, duration time.Duration) {
	// 基础实现，子类可以重写
}

func (b *BaseClient) Request(ctx context.Context, req *APIRequest) (*APIResponse, error) {
	// 构建完整URL
	fullURL, err := b.buildURL(req.Path, req.Params)
	if err != nil {
		return nil, err
	}
	
	// 构建请求体
	var body io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bodyBytes)
	}
	
	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, fullURL, body)
	if err != nil {
		return nil, err
	}
	
	// 设置请求头
	for key, value := range b.headers {
		httpReq.Header.Set(key, value)
	}
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}
	
	// 设置超时
	client := b.client
	if req.Timeout > 0 {
		client = &http.Client{Timeout: req.Timeout}
	}
	
	// 发送请求
	start := time.Now()
	resp, err := client.Do(httpReq)
	duration := time.Since(start)
	
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// 构建响应
	apiResp := &APIResponse{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
		Body:       respBody,
		Duration:   duration,
	}
	
	// 复制响应头
	for key, values := range resp.Header {
		if len(values) > 0 {
			apiResp.Headers[key] = values[0]
		}
	}
	
	return apiResp, nil
}

func (b *BaseClient) buildURL(path string, params map[string]interface{}) (string, error) {
	baseURL, err := url.Parse(b.baseURL)
	if err != nil {
		return "", err
	}
	
	fullURL, err := baseURL.Parse(path)
	if err != nil {
		return "", err
	}
	
	// 添加查询参数
	if len(params) > 0 {
		query := fullURL.Query()
		for key, value := range params {
			query.Set(key, fmt.Sprintf("%v", value))
		}
		fullURL.RawQuery = query.Encode()
	}
	
	return fullURL.String(), nil
}
