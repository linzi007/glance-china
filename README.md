# Glance 中国版

一个专为中国用户优化的自托管仪表板，将所有信息流整合到一个地方。

## ✨ 特性

### 🇨🇳 中国本土化
- **Bilibili 视频**：支持多UP主视频流，实时更新
- **知乎热榜**：热门话题和问答，支持分类筛选
- **Gitee 仓库**：代码仓库状态，Issues 和 PR 跟踪
- **微博热搜**：实时热搜榜单和话题趋势
- **斗鱼直播**：主播状态和热门游戏
- **掘金技术**：技术文章和开发者动态

### 🚀 性能优化
- 多级缓存系统（内存 + Redis + 磁盘）
- 智能限流和并发控制
- 连接池管理和资源优化
- 实时性能监控和指标收集

### 🌏 完整本地化
- 中文界面和时间格式
- 中文数字格式化（万、亿）
- 本地化错误消息和提示
- 支持中英文切换

### 🔧 易于部署
- Docker 一键部署
- 完整的配置示例
- 健康检查和监控
- 平滑升级支持

## 🚀 快速开始

### Docker 部署（推荐）

```bash
# 克隆项目
git clone https://github.com/linzi007/glance-china.git
cd glance-china

# 复制配置文件
cp examples/glance-china.yml glance.yml

# 启动服务
docker-compose up -d
```

### 二进制部署

```bash
# 下载最新版本
wget https://github.com/linzi007/glance-china/releases/latest/download/glance-china-linux-amd64.tar.gz
tar -xzf glance-china-linux-amd64.tar.gz

# 运行
./glance-china --config glance.yml
```

## 📋 配置说明

### 基础配置

```yaml
server:
  host: 0.0.0.0
  port: 8080
  region: cn
  language: zh-CN

# API 配置
api-sources:
  bilibili:
    base-url: https://api.bilibili.com
    rate-limit: 100
  zhihu:
    base-url: https://www.zhihu.com/api
    rate-limit: 50
```

### 组件配置示例

```yaml
widgets:
  # Bilibili 视频
  - type: bilibili-videos
    up-masters:
      - uid: "123456"
        name: "技术UP主"
    limit: 10
    style: grid

  # 知乎热榜
  - type: zhihu-trending
    categories:
      - technology
      - science
    limit: 15
    show-images: true

  # Gitee 仓库
  - type: gitee-repos
    repositories:
      - "owner/repo-name"
    token: ${GITEE_TOKEN}
    show-issues: true
```

## 🔧 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `GLANCE_CONFIG` | 配置文件路径 | `glance.yml` |
| `GLANCE_PORT` | 服务端口 | `8080` |
| `GLANCE_REGION` | 区域设置 | `cn` |
| `REDIS_URL` | Redis 连接地址 | `redis://localhost:6379` |
| `GITEE_TOKEN` | Gitee API Token | - |
| `WEIBO_APP_KEY` | 微博应用密钥 | - |

## 📊 性能监控

访问 `http://localhost:8080/metrics` 查看性能指标：

- 请求响应时间
- 缓存命中率
- API 调用统计
- 内存使用情况
- 错误率统计

## 🔄 从原版 Glance 迁移

使用内置的迁移工具：

```bash
./glance-china migrate --from original-glance.yml --to glance-china.yml
```

自动转换规则：
- `videos` (YouTube) → `bilibili-videos`
- `reddit` → `zhihu-trending`
- `twitch-channels` → `douyu-live`

## 🛠️ 开发

### 构建项目

```bash
# 安装依赖
go mod download

# 构建
make build

# 运行测试
make test

# 性能测试
make benchmark
```

### 添加新组件

1. 在 `internal/widget/` 创建新组件
2. 实现 `Widget` 接口
3. 在 `internal/widget/registry.go` 注册
4. 添加配置示例和文档

## 📝 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

- GitHub Issues: [问题反馈](https://github.com/linzi007/glance-china/issues)