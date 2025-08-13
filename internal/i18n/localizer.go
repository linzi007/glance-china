package i18n

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Localizer 本地化器
type Localizer struct {
	locale   string
	messages map[string]string
	timeZone *time.Location
}

// SupportedLocales 支持的语言
var SupportedLocales = []string{"zh-CN", "en-US"}

// NewLocalizer 创建本地化器
func NewLocalizer(locale string) *Localizer {
	if locale == "" {
		locale = "zh-CN" // 默认中文
	}
	
	// 设置时区
	timeZone, _ := time.LoadLocation("Asia/Shanghai")
	if locale == "en-US" {
		timeZone, _ = time.LoadLocation("UTC")
	}
	
	return &Localizer{
		locale:   locale,
		messages: getMessages(locale),
		timeZone: timeZone,
	}
}

// T 翻译文本
func (l *Localizer) T(key string, args ...interface{}) string {
	message, exists := l.messages[key]
	if !exists {
		return key // 如果没有翻译，返回原始key
	}
	
	if len(args) > 0 {
		return fmt.Sprintf(message, args...)
	}
	
	return message
}

// FormatTime 格式化时间
func (l *Localizer) FormatTime(t time.Time, format string) string {
	localTime := t.In(l.timeZone)
	
	switch l.locale {
	case "zh-CN":
		return l.formatChineseTime(localTime, format)
	default:
		return localTime.Format(format)
	}
}

// FormatRelativeTime 格式化相对时间
func (l *Localizer) FormatRelativeTime(t time.Time) string {
	now := time.Now().In(l.timeZone)
	diff := now.Sub(t.In(l.timeZone))
	
	switch l.locale {
	case "zh-CN":
		return l.formatChineseRelativeTime(diff)
	default:
		return l.formatEnglishRelativeTime(diff)
	}
}

// FormatNumber 格式化数字
func (l *Localizer) FormatNumber(num int64) string {
	switch l.locale {
	case "zh-CN":
		return l.formatChineseNumber(num)
	default:
		return l.formatEnglishNumber(num)
	}
}

// FormatDuration 格式化时长
func (l *Localizer) FormatDuration(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	
	switch l.locale {
	case "zh-CN":
		if hours > 0 {
			return fmt.Sprintf("%d小时%d分钟", hours, minutes)
		} else if minutes > 0 {
			return fmt.Sprintf("%d分%d秒", minutes, secs)
		} else {
			return fmt.Sprintf("%d秒", secs)
		}
	default:
		if hours > 0 {
			return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
		} else {
			return fmt.Sprintf("%d:%02d", minutes, secs)
		}
	}
}

func (l *Localizer) formatChineseTime(t time.Time, format string) string {
	switch format {
	case "date":
		return t.Format("2006年01月02日")
	case "time":
		return t.Format("15:04")
	case "datetime":
		return t.Format("2006年01月02日 15:04")
	case "full":
		weekday := l.getChineseWeekday(t.Weekday())
		return fmt.Sprintf("%s %s", t.Format("2006年01月02日"), weekday)
	default:
		return t.Format(format)
	}
}

func (l *Localizer) formatChineseRelativeTime(diff time.Duration) string {
	if diff < time.Minute {
		return "刚刚"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d分钟前", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		return fmt.Sprintf("%d小时前", hours)
	} else if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d天前", days)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		return fmt.Sprintf("%d个月前", months)
	} else {
		years := int(diff.Hours() / (24 * 365))
		return fmt.Sprintf("%d年前", years)
	}
}

func (l *Localizer) formatEnglishRelativeTime(diff time.Duration) string {
	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 30*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := int(diff.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

func (l *Localizer) formatChineseNumber(num int64) string {
	if num < 1000 {
		return strconv.FormatInt(num, 10)
	} else if num < 10000 {
		return fmt.Sprintf("%.1f千", float64(num)/1000)
	} else if num < 100000000 {
		return fmt.Sprintf("%.1f万", float64(num)/10000)
	} else {
		return fmt.Sprintf("%.1f亿", float64(num)/100000000)
	}
}

func (l *Localizer) formatEnglishNumber(num int64) string {
	if num < 1000 {
		return strconv.FormatInt(num, 10)
	} else if num < 1000000 {
		return fmt.Sprintf("%.1fK", float64(num)/1000)
	} else if num < 1000000000 {
		return fmt.Sprintf("%.1fM", float64(num)/1000000)
	} else {
		return fmt.Sprintf("%.1fB", float64(num)/1000000000)
	}
}

func (l *Localizer) getChineseWeekday(weekday time.Weekday) string {
	weekdays := map[time.Weekday]string{
		time.Sunday:    "星期日",
		time.Monday:    "星期一",
		time.Tuesday:   "星期二",
		time.Wednesday: "星期三",
		time.Thursday:  "星期四",
		time.Friday:    "星期五",
		time.Saturday:  "星期六",
	}
	return weekdays[weekday]
}

// GetLocale 获取当前语言
func (l *Localizer) GetLocale() string {
	return l.locale
}

// IsChineseLocale 是否为中文语言
func (l *Localizer) IsChineseLocale() bool {
	return strings.HasPrefix(l.locale, "zh")
}
