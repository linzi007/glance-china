package widget

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// ZhihuTrendingWidget 知乎热榜组件
type ZhihuTrendingWidget struct {
	ChineseWidget
	Categories    []string `yaml:"categories"`
	Limit         int      `yaml:"limit"`
	ShowImages    bool     `yaml:"show-images"`
	CollapseAfter int      `yaml:"collapse-after"`
}

type ZhihuTrendingData struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	Excerpt             string    `json:"excerpt"`
	URL                 string    `json:"url"`
	Author              string    `json:"author"`
	AuthorURL           string    `json:"author_url"`
	Image               string    `json:"image,omitempty"`
	HeatValue           int64     `json:"heat_value"`
	AnswerCount         int       `json:"answer_count"`
	UpdatedAt           time.Time `json:"updated_at"`
	Category            string    `json:"category"`
	HeatValueFormatted  string    `json:"heat_value_formatted"`
	UpdatedAtFormatted  string    `json:"updated_at_formatted"`
	CategoryLocalized   string    `json:"category_localized"`
}

type ZhihuAPIResponse struct {
	Data []struct {
		Target struct {
			ID       int64  `json:"id"`
			Title    string `json:"title"`
			Excerpt  string `json:"excerpt"`
			URL      string `json:"url"`
			Author   struct {
				Name      string `json:"name"`
				URL       string `json:"url"`
				AvatarURL string `json:"avatar_url"`
			} `json:"author"`
		} `json:"target"`
		DetailText string `json:"detail_text"`
		HeatValue  int64  `json:"heat_value"`
	} `json:"data"`
}

func NewZhihuTrendingWidget() *ZhihuTrendingWidget {
	return &ZhihuTrendingWidget{
		ChineseWidget: ChineseWidget{
			BaseWidget: BaseWidget{
				Type: "zhihu-trending",
			},
			Region:    "cn",
			APISource: "zhihu",
		},
		Limit:         15,
		ShowImages:    true,
		CollapseAfter: 5,
		Categories:    []string{"all"},
	}
}

func (z *ZhihuTrendingWidget) GetData(ctx context.Context, config Config) (interface{}, error) {
	if locale := ctx.Value("locale"); locale != nil {
		z.InitLocalizer(locale.(string))
	}

	trending, err := z.fetchTrending(ctx)
	if err != nil {
		return nil, err
	}

	if len(z.Categories) > 0 && z.Categories[0] != "all" {
		trending = z.filterByCategories(trending)
	}

	if len(trending) > z.Limit {
		trending = trending[:z.Limit]
	}

	localizer := z.GetLocalizer()
	for i := range trending {
		trending[i].HeatValueFormatted = localizer.FormatNumber(trending[i].HeatValue)
		trending[i].UpdatedAtFormatted = localizer.FormatRelativeTime(trending[i].UpdatedAt)
		trending[i].CategoryLocalized = z.localizeCategory(trending[i].Category, localizer)
	}

	return map[string]interface{}{
		"trending":       trending,
		"show_images":    z.ShowImages,
		"collapse_after": z.CollapseAfter,
		"title":          z.getLocalizedTitle(),
		"locale":         localizer.GetLocale(),
		"labels": map[string]string{
			"heat":      localizer.T("热度"),
			"answers":   localizer.T("回答"),
			"updated":   localizer.T("更新时间"),
			"category":  localizer.T("分类"),
		},
	}, nil
}

func (z *ZhihuTrendingWidget) fetchTrending(ctx context.Context) ([]ZhihuTrendingData, error) {
	// 使用知乎热榜API（这里使用模拟的API端点）
	url := "https://www.zhihu.com/api/v3/feed/topstory/hot-lists/total"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://www.zhihu.com")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp ZhihuAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var trending []ZhihuTrendingData
	for _, item := range apiResp.Data {
		trending = append(trending, ZhihuTrendingData{
			ID:          fmt.Sprintf("%d", item.Target.ID),
			Title:       item.Target.Title,
			Excerpt:     item.Target.Excerpt,
			URL:         item.Target.URL,
			Author:      item.Target.Author.Name,
			AuthorURL:   item.Target.Author.URL,
			Image:       item.Target.Author.AvatarURL,
			HeatValue:   item.HeatValue,
			UpdatedAt:   time.Now(),
			Category:    "general",
		})
	}

	return trending, nil
}

func (z *ZhihuTrendingWidget) filterByCategories(trending []ZhihuTrendingData) []ZhihuTrendingData {
	// 实现分类过滤逻辑
	var filtered []ZhihuTrendingData
	categoryMap := make(map[string]bool)
	for _, cat := range z.Categories {
		categoryMap[cat] = true
	}

	for _, item := range trending {
		if categoryMap[item.Category] {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func (z *ZhihuTrendingWidget) GetCacheKey(config Config) string {
	return fmt.Sprintf("zhihu-trending:%v", z.Categories)
}

func (z *ZhihuTrendingWidget) getTitle() string {
	if z.Title != "" {
		return z.Title
	}
	return "知乎热榜"
}

func (z *ZhihuTrendingWidget) Validate(config Config) error {
	if z.Limit <= 0 {
		return fmt.Errorf("limit 必须大于 0")
	}
	return nil
}

func (z *ZhihuTrendingWidget) localizeCategory(category string, localizer *i18n.Localizer) string {
	key := fmt.Sprintf("category.%s", category)
	translated := localizer.T(key)
	if translated == key {
		return category // 如果没有翻译，返回原文
	}
	return translated
}

func (z *ZhihuTrendingWidget) getLocalizedTitle() string {
	if z.Title != "" {
		return z.Title
	}
	return z.T("widget.zhihu_trending")
}
