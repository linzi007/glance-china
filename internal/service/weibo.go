package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WeiboClient 微博客户端
type WeiboClient struct {
	BaseClient
	appKey    string
	appSecret string
}

// WeiboHotSearchData 微博热搜数据
type WeiboHotSearchData struct {
	Rank     int    `json:"rank"`
	Keyword  string `json:"keyword"`
	URL      string `json:"url"`
	HotValue int64  `json:"hot_value"`
	Category string `json:"category"`
	Icon     string `json:"icon,omitempty"`
}

type WeiboAPIResponse struct {
	Data []struct {
		Realpos   int    `json:"realpos"`
		Word      string `json:"word"`
		WordScheme string `json:"word_scheme"`
		Num       int64  `json:"num"`
		Category  string `json:"category"`
		Icon      string `json:"icon"`
	} `json:"data"`
}

func NewWeiboClient(config APISourceConfig) *WeiboClient {
	return &WeiboClient{
		BaseClient: BaseClient{
			name:      "weibo",
			baseURL:   config.BaseURL,
			timeout:   config.Timeout,
			headers:   config.Headers,
			client:    &http.Client{Timeout: config.Timeout},
		},
		appKey:    config.Headers["app-key"],
		appSecret: config.Token,
	}
}

// GetHotSearch 获取微博热搜
func (w *WeiboClient) GetHotSearch(ctx context.Context) ([]WeiboHotSearchData, error) {
	req := &APIRequest{
		Method:  "GET",
		Path:    "/2/search/topics.json",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		},
		Timeout: 10 * time.Second,
	}
	
	resp, err := w.Request(ctx, req)
	if err != nil {
		return nil, err
	}
	
	var apiResp WeiboAPIResponse
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, err
	}
	
	var hotSearches []WeiboHotSearchData
	for i, item := range apiResp.Data {
		if i >= 50 { // 限制数量
			break
		}
		
		hotSearches = append(hotSearches, WeiboHotSearchData{
			Rank:     item.Realpos,
			Keyword:  item.Word,
			URL:      fmt.Sprintf("https://s.weibo.com/weibo?q=%s", item.WordScheme),
			HotValue: item.Num,
			Category: item.Category,
			Icon:     item.Icon,
		})
	}
	
	return hotSearches, nil
}
