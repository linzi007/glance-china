package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DouyuClient 斗鱼客户端
type DouyuClient struct {
	BaseClient
}

// DouyuStreamData 斗鱼直播数据
type DouyuStreamData struct {
	RoomID      string `json:"room_id"`
	RoomName    string `json:"room_name"`
	OwnerName   string `json:"owner_name"`
	OwnerAvatar string `json:"owner_avatar"`
	GameName    string `json:"game_name"`
	Viewers     int    `json:"viewers"`
	IsLive      bool   `json:"is_live"`
	StreamURL   string `json:"stream_url"`
	Thumbnail   string `json:"thumbnail"`
}

type DouyuAPIResponse struct {
	Error int `json:"error"`
	Data  []struct {
		RoomID     string `json:"room_id"`
		RoomName   string `json:"room_name"`
		OwnerName  string `json:"owner_name"`
		Avatar     string `json:"avatar"`
		GameName   string `json:"game_name"`
		Online     int    `json:"online"`
		ShowStatus int    `json:"show_status"`
		RoomSrc    string `json:"room_src"`
	} `json:"data"`
}

func NewDouyuClient(config APISourceConfig) *DouyuClient {
	return &DouyuClient{
		BaseClient: BaseClient{
			name:      "douyu",
			baseURL:   config.BaseURL,
			timeout:   config.Timeout,
			headers:   config.Headers,
			client:    &http.Client{Timeout: config.Timeout},
		},
	}
}

// GetLiveStreams 获取直播流信息
func (d *DouyuClient) GetLiveStreams(ctx context.Context, roomIDs []string) ([]DouyuStreamData, error) {
	var allStreams []DouyuStreamData
	
	for _, roomID := range roomIDs {
		stream, err := d.getRoomInfo(ctx, roomID)
		if err != nil {
			continue // 跳过错误的房间
		}
		allStreams = append(allStreams, *stream)
	}
	
	return allStreams, nil
}

func (d *DouyuClient) getRoomInfo(ctx context.Context, roomID string) (*DouyuStreamData, error) {
	req := &APIRequest{
		Method: "GET",
		Path:   fmt.Sprintf("/api/v1/room/%s", roomID),
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		},
		Timeout: 5 * time.Second,
	}
	
	resp, err := d.Request(ctx, req)
	if err != nil {
		return nil, err
	}
	
	var apiResp DouyuAPIResponse
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, err
	}
	
	if apiResp.Error != 0 || len(apiResp.Data) == 0 {
		return nil, fmt.Errorf("room not found: %s", roomID)
	}
	
	room := apiResp.Data[0]
	return &DouyuStreamData{
		RoomID:      room.RoomID,
		RoomName:    room.RoomName,
		OwnerName:   room.OwnerName,
		OwnerAvatar: room.Avatar,
		GameName:    room.GameName,
		Viewers:     room.Online,
		IsLive:      room.ShowStatus == 1,
		StreamURL:   fmt.Sprintf("https://www.douyu.com/%s", room.RoomID),
		Thumbnail:   room.RoomSrc,
	}, nil
}
