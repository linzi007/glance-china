package i18n

// getMessages 获取指定语言的翻译消息
func getMessages(locale string) map[string]string {
	switch locale {
	case "zh-CN":
		return chineseMessages
	case "en-US":
		return englishMessages
	default:
		return chineseMessages
	}
}

// chineseMessages 中文翻译
var chineseMessages = map[string]string{
	// 通用
	"loading":           "加载中...",
	"error":             "错误",
	"retry":             "重试",
	"show_more":         "显示更多",
	"show_less":         "收起",
	"refresh":           "刷新",
	"settings":          "设置",
	"about":             "关于",
	
	// 时间相关
	"time.just_now":     "刚刚",
	"time.minutes_ago":  "%d分钟前",
	"time.hours_ago":    "%d小时前",
	"time.days_ago":     "%d天前",
	"time.months_ago":   "%d个月前",
	"time.years_ago":    "%d年前",
	"time.today":        "今天",
	"time.yesterday":    "昨天",
	"time.tomorrow":     "明天",
	
	// 数字单位
	"number.thousand":   "千",
	"number.ten_thousand": "万",
	"number.hundred_million": "亿",
	"number.views":      "播放",
	"number.likes":      "点赞",
	"number.comments":   "评论",
	"number.shares":     "分享",
	"number.followers":  "粉丝",
	"number.stars":      "星标",
	"number.forks":      "分叉",
	"number.issues":     "问题",
	"number.viewers":    "观看人数",
	
	// 组件标题
	"widget.bilibili_videos":    "Bilibili 视频",
	"widget.zhihu_trending":     "知乎热榜",
	"widget.gitee_repos":        "Gitee 仓库",
	"widget.weibo_hot_search":   "微博热搜",
	"widget.douyu_live":         "斗鱼直播",
	"widget.weather":            "天气",
	"widget.calendar":           "日历",
	"widget.clock":              "时钟",
	"widget.rss":                "RSS 订阅",
	"widget.bookmarks":          "书签",
	"widget.server_stats":       "服务器状态",
	"widget.docker_containers":  "Docker 容器",
	
	// 状态
	"status.online":     "在线",
	"status.offline":    "离线",
	"status.live":       "直播中",
	"status.not_live":   "未直播",
	"status.healthy":    "正常",
	"status.unhealthy":  "异常",
	"status.running":    "运行中",
	"status.stopped":    "已停止",
	"status.error":      "错误",
	
	// 操作
	"action.view":       "查看",
	"action.watch":      "观看",
	"action.read":       "阅读",
	"action.download":   "下载",
	"action.share":      "分享",
	"action.like":       "点赞",
	"action.comment":    "评论",
	"action.follow":     "关注",
	"action.star":       "收藏",
	"action.fork":       "分叉",
	
	// 分类
	"category.technology":   "科技",
	"category.entertainment": "娱乐",
	"category.gaming":       "游戏",
	"category.music":        "音乐",
	"category.sports":       "体育",
	"category.news":         "新闻",
	"category.education":    "教育",
	"category.lifestyle":    "生活",
	"category.travel":       "旅行",
	"category.food":         "美食",
	
	// 错误消息
	"error.network":         "网络连接错误",
	"error.timeout":         "请求超时",
	"error.rate_limit":      "请求频率过高",
	"error.not_found":       "未找到内容",
	"error.server_error":    "服务器错误",
	"error.invalid_config":  "配置错误",
	"error.auth_failed":     "认证失败",
	
	// 配置
	"config.theme":          "主题",
	"config.language":       "语言",
	"config.timezone":       "时区",
	"config.refresh_rate":   "刷新频率",
	"config.cache_duration": "缓存时长",
	
	// 天气
	"weather.sunny":         "晴天",
	"weather.cloudy":        "多云",
	"weather.rainy":         "雨天",
	"weather.snowy":         "雪天",
	"weather.foggy":         "雾天",
	"weather.windy":         "大风",
	"weather.temperature":   "温度",
	"weather.humidity":      "湿度",
	"weather.pressure":      "气压",
	"weather.visibility":    "能见度",
	
	// 直播
	"live.viewers":          "观看人数",
	"live.duration":         "直播时长",
	"live.category":         "分类",
	"live.title":            "标题",
	"live.streamer":         "主播",
	
	// 仓库
	"repo.stars":            "星标数",
	"repo.forks":            "分叉数",
	"repo.issues":           "问题数",
	"repo.pull_requests":    "拉取请求",
	"repo.last_commit":      "最后提交",
	"repo.language":         "编程语言",
	"repo.license":          "许可证",
	"repo.size":             "大小",
	
	// 视频
	"video.duration":        "时长",
	"video.views":           "播放量",
	"video.likes":           "点赞数",
	"video.comments":        "评论数",
	"video.published":       "发布时间",
	"video.author":          "作者",
	"video.channel":         "频道",
}

// englishMessages 英文翻译
var englishMessages = map[string]string{
	// 通用
	"loading":           "Loading...",
	"error":             "Error",
	"retry":             "Retry",
	"show_more":         "Show More",
	"show_less":         "Show Less",
	"refresh":           "Refresh",
	"settings":          "Settings",
	"about":             "About",
	
	// 时间相关
	"time.just_now":     "just now",
	"time.minutes_ago":  "%d minutes ago",
	"time.hours_ago":    "%d hours ago",
	"time.days_ago":     "%d days ago",
	"time.months_ago":   "%d months ago",
	"time.years_ago":    "%d years ago",
	"time.today":        "today",
	"time.yesterday":    "yesterday",
	"time.tomorrow":     "tomorrow",
	
	// 数字单位
	"number.thousand":   "K",
	"number.ten_thousand": "K",
	"number.hundred_million": "M",
	"number.views":      "views",
	"number.likes":      "likes",
	"number.comments":   "comments",
	"number.shares":     "shares",
	"number.followers":  "followers",
	"number.stars":      "stars",
	"number.forks":      "forks",
	"number.issues":     "issues",
	"number.viewers":    "viewers",
	
	// 组件标题
	"widget.bilibili_videos":    "Bilibili Videos",
	"widget.zhihu_trending":     "Zhihu Trending",
	"widget.gitee_repos":        "Gitee Repositories",
	"widget.weibo_hot_search":   "Weibo Hot Search",
	"widget.douyu_live":         "Douyu Live",
	"widget.weather":            "Weather",
	"widget.calendar":           "Calendar",
	"widget.clock":              "Clock",
	"widget.rss":                "RSS Feeds",
	"widget.bookmarks":          "Bookmarks",
	"widget.server_stats":       "Server Stats",
	"widget.docker_containers":  "Docker Containers",
	
	// 其他英文翻译...
	"status.online":     "Online",
	"status.offline":    "Offline",
	"status.live":       "Live",
	"status.not_live":   "Not Live",
	"status.healthy":    "Healthy",
	"status.unhealthy":  "Unhealthy",
	"status.running":    "Running",
	"status.stopped":    "Stopped",
	"status.error":      "Error",
}
