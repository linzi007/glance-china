# 故障排除指南

## 常见问题

### 1. 服务启动失败

#### 问题：端口被占用
```
Error: listen tcp :8080: bind: address already in use
```

**解决方案：**
```bash
# 查找占用端口的进程
sudo netstat -tlnp | grep :8080
sudo lsof -i :8080

# 杀死进程或更改端口
sudo kill -9 <PID>
# 或在配置文件中修改端口
```

#### 问题：配置文件错误
```
Error: yaml: unmarshal errors
```

**解决方案：**
```bash
# 验证配置文件语法
./glance-china --config glance.yml --validate

# 检查 YAML 格式
yamllint glance.yml
```

### 2. API 调用失败

#### 问题：Bilibili API 无响应
```
Error: failed to fetch bilibili videos: context deadline exceeded
```

**解决方案：**
```bash
# 测试网络连接
curl -I https://api.bilibili.com

# 检查 DNS 解析
nslookup api.bilibili.com

# 增加超时时间
api-sources:
  bilibili:
    timeout: 60s
    retries: 5
```

#### 问题：API 限流
```
Error: rate limit exceeded for bilibili API
```

**解决方案：**
```yaml
# 调整限流配置
api-sources:
  bilibili:
    rate-limit: 50  # 降低请求频率
    
# 增加缓存时间
- type: bilibili-videos
  cache-duration: 1800s  # 30分钟
```

### 3. 内存使用过高

#### 问题：内存泄漏
```
Error: runtime: out of memory
```

**解决方案：**
```bash
# 查看内存使用情况
curl http://localhost:8080/debug/pprof/heap

# 调整缓存大小
performance:
  cache:
    memory-size: 50MB  # 减少内存缓存
    
# 启用垃圾回收调试
export GODEBUG=gctrace=1
```

### 4. 缓存问题

#### 问题：Redis 连接失败
```
Error: failed to connect to redis: dial tcp: connection refused
```

**解决方案：**
```bash
# 检查 Redis 服务
sudo systemctl status redis
redis-cli ping

# 检查连接配置
export REDIS_URL="redis://localhost:6379"

# 禁用 Redis 缓存（临时）
cache:
  levels:
    - type: memory
      size: 200MB
```

### 5. 组件显示异常

#### 问题：组件数据为空
```
Widget shows "No data available"
```

**解决方案：**
```bash
# 检查 API 响应
curl "https://api.bilibili.com/x/space/arc/search?mid=123456"

# 验证配置参数
- type: bilibili-videos
  up-masters:
    - uid: "123456"  # 确保 UID 正确
      name: "UP主名称"

# 查看详细日志
logging:
  level: debug
```

#### 问题：中文显示乱码
```
Chinese characters show as "???"
```

**解决方案：**
```yaml
# 设置正确的语言
server:
  language: zh-CN
  
# 检查系统编码
export LANG=zh_CN.UTF-8
export LC_ALL=zh_CN.UTF-8
```

## 性能问题

### 1. 响应速度慢

**诊断步骤：**
```bash
# 查看性能指标
curl http://localhost:8080/metrics

# 检查网络延迟
ping api.bilibili.com
traceroute api.bilibili.com

# 分析慢查询
curl http://localhost:8080/debug/requests
```

**优化方案：**
```yaml
# 启用并发请求
performance:
  workers: 20
  
# 优化缓存策略
cache:
  levels:
    - type: memory
      size: 200MB
      ttl: 600s
      
# 减少刷新频率
widgets:
  - type: bilibili-videos
    refresh-interval: 600s  # 10分钟
```

### 2. CPU 使用率高

**诊断步骤：**
```bash
# CPU 性能分析
curl http://localhost:8080/debug/pprof/profile > cpu.prof
go tool pprof cpu.prof

# 查看协程状态
curl http://localhost:8080/debug/pprof/goroutine
```

**优化方案：**
```yaml
# 限制并发数
performance:
  workers: 5
  
# 增加请求间隔
rate-limit:
  global:
    requests-per-minute: 300
```

## 网络问题

### 1. DNS 解析失败

```bash
# 测试 DNS
nslookup api.bilibili.com
dig api.bilibili.com

# 使用公共 DNS
echo "nameserver 8.8.8.8" >> /etc/resolv.conf
```

### 2. 防火墙阻拦

```bash
# 检查防火墙规则
sudo iptables -L
sudo ufw status

# 开放端口
sudo ufw allow 8080
sudo iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
```

### 3. 代理配置

```bash
# 设置代理
export HTTP_PROXY=http://proxy.company.com:8080
export HTTPS_PROXY=http://proxy.company.com:8080

# 在配置中设置代理
api-sources:
  bilibili:
    proxy: http://proxy.company.com:8080
```

## 日志分析

### 启用详细日志

```yaml
logging:
  level: debug
  format: json
  output: /var/log/glance-china.log
```

### 常见错误日志

#### API 超时
```json
{
  "level": "error",
  "msg": "API request timeout",
  "widget": "bilibili-videos",
  "url": "https://api.bilibili.com/x/space/arc/search",
  "timeout": "30s"
}
```

#### 缓存失效
```json
{
  "level": "warn",
  "msg": "Cache miss",
  "key": "bilibili:videos:123456",
  "reason": "expired"
}
```

#### 限流触发
```json
{
  "level": "warn",
  "msg": "Rate limit exceeded",
  "widget": "zhihu-trending",
  "limit": "50/min",
  "current": "51"
}
```

## 监控和告警

### 健康检查

```bash
# 基础健康检查
curl http://localhost:8080/health

# 详细健康状态
curl http://localhost:8080/health?detailed=true
```

### 指标监控

```bash
# Prometheus 格式指标
curl http://localhost:8080/metrics

# 关键指标
# - glance_requests_total
# - glance_request_duration_seconds
# - glance_cache_hit_ratio
# - glance_api_errors_total
```

### 告警规则示例

```yaml
# Prometheus 告警规则
groups:
  - name: glance-china
    rules:
      - alert: GlanceHighErrorRate
        expr: rate(glance_api_errors_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Glance API error rate is high"
          
      - alert: GlanceCacheLowHitRate
        expr: glance_cache_hit_ratio < 0.7
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Glance cache hit rate is low"
```

## 数据恢复

### 缓存数据恢复

```bash
# 清理损坏的缓存
rm -rf ./cache/*
redis-cli FLUSHALL

# 重新加载数据
curl -X POST http://localhost:8080/api/refresh
```

### 配置备份恢复

```bash
# 备份配置
cp glance.yml glance.yml.backup.$(date +%Y%m%d)

# 恢复配置
cp glance.yml.backup.20240101 glance.yml
sudo systemctl restart glance-china
```

## 获取帮助

如果以上方法都无法解决问题，请：

1. 收集详细日志和错误信息
2. 记录复现步骤
3. 提供系统环境信息
4. 在 GitHub 提交 Issue：https://github.com/linzi007/glance-china/issues

### 问题报告模板

```markdown
## 问题描述
简要描述遇到的问题

## 环境信息
- 操作系统：
- Glance 版本：
- Go 版本：
- 部署方式：

## 复现步骤
1. 
2. 
3. 

## 期望结果
描述期望的正常行为

## 实际结果
描述实际发生的情况

## 日志信息
```
相关的错误日志
```

## 配置文件
```yaml
# 相关的配置片段
