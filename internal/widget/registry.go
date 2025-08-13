package widget

import (
	"fmt"
	"sync"
)

// Registry 组件注册表
type Registry struct {
	widgets map[string]func() Widget
	mu      sync.RWMutex
}

var globalRegistry = &Registry{
	widgets: make(map[string]func() Widget),
}

// RegisterWidget 注册组件
func RegisterWidget(widgetType string, constructor func() Widget) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.widgets[widgetType] = constructor
}

// CreateWidget 创建组件实例
func CreateWidget(widgetType string) (Widget, error) {
	globalRegistry.mu.RLock()
	constructor, exists := globalRegistry.widgets[widgetType]
	globalRegistry.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("unknown widget type: %s", widgetType)
	}
	
	return constructor(), nil
}

// GetRegisteredWidgets 获取所有已注册的组件类型
func GetRegisteredWidgets() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	
	var types []string
	for widgetType := range globalRegistry.widgets {
		types = append(types, widgetType)
	}
	return types
}

// 初始化注册中国版组件
func init() {
	RegisterWidget("bilibili-videos", func() Widget {
		return NewBilibiliVideosWidget()
	})
	
	RegisterWidget("zhihu-trending", func() Widget {
		return NewZhihuTrendingWidget()
	})
	
	RegisterWidget("gitee-repos", func() Widget {
		return NewGiteeReposWidget()
	})
}
