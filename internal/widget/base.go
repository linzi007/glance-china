package widget

import (
	"context"
	"time"
	
	"github.com/glance-china/internal/i18n"
)

// Widget 定义所有组件的基础接口
type Widget interface {
	GetType() string
	GetData(ctx context.Context, config Config) (interface{}, error)
	GetCacheKey(config Config) string
	GetCacheDuration() time.Duration
	Validate(config Config) error
}

// BaseWidget 提供基础实现
type BaseWidget struct {
	Type         string        `yaml:"type"`
	Title        string        `yaml:"title,omitempty"`
	TitleURL     string        `yaml:"title-url,omitempty"`
	Cache        string        `yaml:"cache,omitempty"`
	CSSClass     string        `yaml:"css-class,omitempty"`
	Locale       string        `yaml:"locale,omitempty"`
	cacheDuration time.Duration
	localizer    *i18n.Localizer
}

func (b *BaseWidget) InitLocalizer(locale string) {
	if locale == "" {
		locale = "zh-CN" // 默认中文
	}
	b.Locale = locale
	b.localizer = i18n.NewLocalizer(locale)
}

func (b *BaseWidget) GetLocalizer() *i18n.Localizer {
	if b.localizer == nil {
		b.InitLocalizer("zh-CN")
	}
	return b.localizer
}

func (b *BaseWidget) T(key string, args ...interface{}) string {
	return b.GetLocalizer().T(key, args...)
}

func (b *BaseWidget) GetType() string {
	return b.Type
}

func (b *BaseWidget) GetCacheDuration() time.Duration {
	if b.cacheDuration == 0 {
		return 5 * time.Minute // 默认缓存5分钟
	}
	return b.cacheDuration
}

func (b *BaseWidget) Validate(config Config) error {
	// 基础验证逻辑
	return nil
}

// ChineseWidget 中国版组件基础结构
type ChineseWidget struct {
	BaseWidget
	Region    string   `yaml:"region,omitempty"`
	APISource string   `yaml:"api-source,omitempty"`
	Fallback  []string `yaml:"fallback,omitempty"`
}

// Config 通用配置接口
type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetStringSlice(key string) []string
}
