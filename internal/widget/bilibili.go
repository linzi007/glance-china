package widget

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// BilibiliVideosWidget Bilibili视频组件
type BilibiliVideosWidget struct {
	ChineseWidget
	UPMasters     []BilibiliUPMaster `yaml:"up-masters"`
	Limit         int                `yaml:"limit"`
	Style         string             `yaml:"style"`
	CollapseAfter int                `yaml:"collapse-after"`
}

type BilibiliUPMaster struct {
	UID  string `yaml:"uid"`
	Name string `yaml:"name"`
}

type BilibiliVideoData struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	AuthorURL   string    `json:"author_url"`
	VideoURL    string    `json:"video_url"`
	Thumbnail   string    `json:"thumbnail"`
	Duration    string    `json:"duration"`
	ViewCount   int64     `json:"view_count"`
	PublishedAt time.Time `json:"published_at"`
	Source      string    `json:"source"`
	ViewCountFormatted   string `json:"view_count_formatted"`
	PublishedAtFormatted string `json:"published_at_formatted"`
	DurationFormatted    string `json:"duration_formatted"`
}

type BilibiliAPIResponse struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Data    struct {
		List struct {
			Vlist []struct {
				Aid         int64  `json:"aid"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Author      string `json:"author"`
				Mid         int64  `json:"mid"`
				Created     int64  `json:"created"`
				Length      string `json:"length"`
				VideoReview int64  `json:"video_review"`
				Pic         string `json:"pic"`
			} `json:"vlist"`
		} `json:"list"`
	} `json:"data"`
}

func NewBilibiliVideosWidget() *BilibiliVideosWidget {
	return &BilibiliVideosWidget{
		ChineseWidget: ChineseWidget{
			BaseWidget: BaseWidget{
				Type: "bilibili-videos",
			},
			Region:    "cn",
			APISource: "bilibili",
		},
		Limit:         10,
		Style:         "horizontal-cards",
		CollapseAfter: 5,
	}
}

func (b *BilibiliVideosWidget) GetData(ctx context.Context, config Config) (interface{}, error) {
	if locale := ctx.Value("locale"); locale != nil {
		b.InitLocalizer(locale.(string))
	}
	
	var allVideos []BilibiliVideoData
	
	for _, upMaster := range b.UPMasters {
		videos, err := b.fetchUPMasterVideos(ctx, upMaster)
		if err != nil {
			continue
		}
		allVideos = append(allVideos, videos...)
	}
	
	if len(allVideos) > b.Limit {
		allVideos = allVideos[:b.Limit]
	}
	
	localizer := b.GetLocalizer()
	for i := range allVideos {
		allVideos[i].ViewCountFormatted = localizer.FormatNumber(allVideos[i].ViewCount)
		allVideos[i].PublishedAtFormatted = localizer.FormatRelativeTime(allVideos[i].PublishedAt)
		allVideos[i].DurationFormatted = b.formatDuration(allVideos[i].Duration, localizer)
	}
	
	return map[string]interface{}{
		"videos":        allVideos,
		"style":         b.Style,
		"collapse_after": b.CollapseAfter,
		"title":         b.getLocalizedTitle(),
		"locale":        localizer.GetLocale(),
		"labels": map[string]string{
			"views":     localizer.T("number.views"),
			"published": localizer.T("video.published"),
			"duration":  localizer.T("video.duration"),
			"author":    localizer.T("video.author"),
		},
	}, nil
}

func (b *BilibiliVideosWidget) fetchUPMasterVideos(ctx context.Context, upMaster BilibiliUPMaster) ([]BilibiliVideoData, error) {
	url := fmt.Sprintf("https://api.bilibili.com/x/space/arc/search?mid=%s&ps=10&tid=0&pn=1&keyword=&order=pubdate", upMaster.UID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// 设置必要的请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://www.bilibili.com")
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var apiResp BilibiliAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	
	if apiResp.Code != 0 {
		return nil, fmt.Errorf("bilibili API error: %s", apiResp.Message)
	}
	
	var videos []BilibiliVideoData
	for _, video := range apiResp.Data.List.Vlist {
		videos = append(videos, BilibiliVideoData{
			ID:          strconv.FormatInt(video.Aid, 10),
			Title:       video.Title,
			Author:      upMaster.Name,
			AuthorURL:   fmt.Sprintf("https://space.bilibili.com/%s", upMaster.UID),
			VideoURL:    fmt.Sprintf("https://www.bilibili.com/video/av%d", video.Aid),
			Thumbnail:   video.Pic,
			Duration:    video.Length,
			ViewCount:   video.VideoReview,
			PublishedAt: time.Unix(video.Created, 0),
			Source:      "bilibili",
		})
	}
	
	return videos, nil
}

func (b *BilibiliVideosWidget) GetCacheKey(config Config) string {
	return fmt.Sprintf("bilibili-videos:%d", len(b.UPMasters))
}

func (b *BilibiliVideosWidget) getTitle() string {
	if b.Title != "" {
		return b.Title
	}
	return "Bilibili 视频"
}

func (b *BilibiliVideosWidget) Validate(config Config) error {
	if len(b.UPMasters) == 0 {
		return fmt.Errorf("至少需要配置一个UP主")
	}
	
	for _, upMaster := range b.UPMasters {
		if upMaster.UID == "" {
			return fmt.Errorf("UP主UID不能为空")
		}
	}
	
	return nil
}

func (b *BilibiliVideosWidget) getLocalizedTitle() string {
	if b.Title != "" {
		return b.Title
	}
	return b.T("widget.bilibili_videos")
}

func (b *BilibiliVideosWidget) formatDuration(duration string, localizer *i18n.Localizer) string {
	// 解析时长字符串 "mm:ss" 或 "hh:mm:ss"
	parts := strings.Split(duration, ":")
	if len(parts) == 2 {
		minutes, _ := strconv.Atoi(parts[0])
		seconds, _ := strconv.Atoi(parts[1])
		return localizer.FormatDuration(minutes*60 + seconds)
	} else if len(parts) == 3 {
		hours, _ := strconv.Atoi(parts[0])
		minutes, _ := strconv.Atoi(parts[1])
		seconds, _ := strconv.Atoi(parts[2])
		return localizer.FormatDuration(hours*3600 + minutes*60 + seconds)
	}
	return duration
}
