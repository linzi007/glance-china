package widget

import (
	"context"
	"fmt"
	"time"
	
	"github.com/glance-china/internal/service"
)

// WeiboHotSearchWidget 微博热搜组件
type WeiboHotSearchWidget struct {
	ChineseWidget
	Categories    []string `yaml:"categories"`
	Limit         int      `yaml:"limit"`
	ShowIcons     bool     `yaml:"show-icons"`
	CollapseAfter int      `yaml:"collapse-after"`
}

func NewWeiboHotSearchWidget() *WeiboHotSearchWidget {
	return &WeiboHotSearchWidget{
		ChineseWidget: ChineseWidget{
			BaseWidget: BaseWidget{
				Type: "weibo-hot-search",
			},
			Region:    "cn",
			APISource: "weibo",
		},
		Limit:         20,
		ShowIcons:     true,
		CollapseAfter: 10,
		Categories:    []string{"all"},
	}
}

func (w *WeiboHotSearchWidget) GetData(ctx context.Context, config Config) (interface{}, error) {
	// 获取服务管理器
	serviceManager := ctx.Value("serviceManager").(*service.ServiceManager)
	
	client, err := serviceManager.GetClient("weibo")
	if err != nil {
		return nil, err
	}
	
	weiboClient := client.(*service.WeiboClient)
	hotSearches, err := weiboClient.GetHotSearch(ctx)
	if err != nil {
		return nil, err
	}
	
	// 过滤和限制数量
	if len(hotSearches) > w.Limit {
		hotSearches = hotSearches[:w.Limit]
	}
	
	return map[string]interface{}{
		"hot_searches":   hotSearches,
		"show_icons":     w.ShowIcons,
		"collapse_after": w.CollapseAfter,
		"title":          w.getTitle(),
	}, nil
}

func (w *WeiboHotSearchWidget) GetCacheKey(config Config) string {
	return fmt.Sprintf("weibo-hot-search:%v", w.Categories)
}

func (w *WeiboHotSearchWidget) getTitle() string {
	if w.Title != "" {
		return w.Title
	}
	return "微博热搜"
}

func (w *WeiboHotSearchWidget) Validate(config Config) error {
	if w.Limit <= 0 {
		return fmt.Errorf("limit 必须大于 0")
	}
	return nil
}
