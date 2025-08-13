package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BilibiliClient Bilibili客户端
type BilibiliClient struct {
	BaseClient
}

func NewBilibiliClient(config APISourceConfig) *BilibiliClient {
	return &BilibiliClient{
		BaseClient: BaseClient{
			name:      "bilibili",
			baseURL:   config.BaseURL,
			timeout:   config.Timeout,
			headers:   config.Headers,
			client:    &http.Client{Timeout: config.Timeout},
		},
	}
}

// GetUPMasterVideos 获取UP主视频
func (b *BilibiliClient) GetUPMasterVideos(ctx context.Context, uid string, limit int) ([]BilibiliVideoInfo, error) {
	req := &APIRequest{
		Method: "GET",
		Path:   "/x/space/arc/search",
		Params: map[string]interface{}{
			"mid":     uid,
			"ps":      limit,
			"tid":     0,
			"pn":      1,
			"keyword": "",
			"order":   "pubdate",
		},
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    "https://www.bilibili.com",
		},
		Timeout: 10 * time.Second,
	}
	
	resp, err := b.Request(ctx, req)
	if err != nil {
		return nil, err
	}
	
	var apiResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			List struct {
				Vlist []BilibiliVideoInfo `json:"vlist"`
			} `json:"list"`
		} `json:"data"`
	}
	
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, err
	}
	
	if apiResp.Code != 0 {
		return nil, fmt.Errorf("bilibili API error: %s", apiResp.Message)
	}
	
	return apiResp.Data.List.Vlist, nil
}

type BilibiliVideoInfo struct {
	Aid         int64  `json:"aid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Mid         int64  `json:"mid"`
	Created     int64  `json:"created"`
	Length      string `json:"length"`
	VideoReview int64  `json:"video_review"`
	Pic         string `json:"pic"`
}
