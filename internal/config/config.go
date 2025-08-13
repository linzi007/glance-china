package config

import (
	"fmt"
	"os"
	"time"
	
	"gopkg.in/yaml.v3"
)

// AppConfig 应用配置
type AppConfig struct {
	Server     ServerConfig              `yaml:"server"`
	Theme      ThemeConfig               `yaml:"theme"`
	Pages      []PageConfig              `yaml:"pages"`
	Locale     LocaleConfig              `yaml:"locale"`
}

type LocaleConfig struct {
	Default    string            `yaml:"default"`
	Supported  []string          `yaml:"supported"`
	TimeZone   string            `yaml:"timezone"`
	DateFormat string            `yaml:"date-format"`
	TimeFormat string            `yaml:"time-format"`
	Currency   string            `yaml:"currency"`
	NumberFormat NumberFormatConfig `yaml:"number-format"`
}

type NumberFormatConfig struct {
	DecimalSeparator  string `yaml:"decimal-separator"`
	ThousandSeparator string `yaml:"thousand-separator"`
	UseChineseUnits   bool   `yaml:"use-chinese-units"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port       int                       `yaml:"port"`
	Host       string                    `yaml:"host"`
	Region     string                    `yaml:"region"`
	APISources map[string]APISourceConfig `yaml:"api-sources"`
	Cache      CacheConfig               `yaml:"cache"`
	RateLimit  RateLimitConfig           `yaml:"rate-limit"`
}


// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*AppConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	if config.Locale.Default == "" {
		config.Locale.Default = "zh-CN"
	}
	if len(config.Locale.Supported) == 0 {
		config.Locale.Supported = []string{"zh-CN", "en-US"}
	}
	if config.Locale.TimeZone == "" {
		config.Locale.TimeZone = "Asia/Shanghai"
	}
	if config.Locale.NumberFormat.UseChineseUnits == false && config.Locale.Default == "zh-CN" {
		config.Locale.NumberFormat.UseChineseUnits = true
	}
	
	return &config, nil
}

// ValidateConfig 验证配置
func ValidateConfig(config *AppConfig) error {
	if config.Locale.Default != "" {
		found := false
		for _, supported := range config.Locale.Supported {
			if supported == config.Locale.Default {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("default locale %s is not in supported locales", config.Locale.Default)
		}
	}
	
	// 验证时区
	if config.Locale.TimeZone != "" {
		if _, err := time.LoadLocation(config.Locale.TimeZone); err != nil {
			return fmt.Errorf("invalid timezone: %s", config.Locale.TimeZone)
		}
	}
	
	return nil
}
