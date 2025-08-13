# 部署指南

## 系统要求

### 最低配置
- CPU: 1 核心
- 内存: 512MB
- 存储: 1GB
- 网络: 稳定的互联网连接

### 推荐配置
- CPU: 2 核心
- 内存: 2GB
- 存储: 5GB SSD
- 网络: 100Mbps+

## 部署方式

### 1. Docker 部署（推荐）

#### 单容器部署

\`\`\`bash
# 创建配置目录
mkdir -p /opt/glance-china/config

# 创建配置文件
cat > /opt/glance-china/config/glance.yml << EOF
server:
  host: 0.0.0.0
  port: 8080
  region: cn
  language: zh-CN

widgets:
  - type: bilibili-videos
    up-masters:
      - uid: "123456"
        name: "示例UP主"
    limit: 10
EOF

# 运行容器
docker run -d \
  --name glance-china \
  --restart unless-stopped \
  -p 8080:8080 \
  -v /opt/glance-china/config:/app/config \
  glance-china:latest
\`\`\`

#### Docker Compose 部署

\`\`\`yaml
# docker-compose.yml
version: '3.8'

services:
  glance-china:
    image: glance-china:latest
    container_name: glance-china
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./config:/app/config
      - ./data:/app/data
    environment:
      - GLANCE_CONFIG=/app/config/glance.yml
      - REDIS_URL=redis://redis:6379
    depends_on:
      - redis

  redis:
    image: redis:7-alpine
    container_name: glance-redis
    restart: unless-stopped
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

volumes:
  redis_data:
\`\`\`

### 2. 二进制部署

#### 下载和安装

\`\`\`bash
# 创建用户和目录
sudo useradd -r -s /bin/false glance
sudo mkdir -p /opt/glance-china/{bin,config,data,logs}
sudo chown -R glance:glance /opt/glance-china

# 下载二进制文件
cd /opt/glance-china/bin
sudo wget https://github.com/your-org/glance-china/releases/latest/download/glance-china-linux-amd64.tar.gz
sudo tar -xzf glance-china-linux-amd64.tar.gz
sudo chown glance:glance glance-china
sudo chmod +x glance-china
\`\`\`

#### 创建 systemd 服务

\`\`\`bash
# 创建服务文件
sudo cat > /etc/systemd/system/glance-china.service << EOF
[Unit]
Description=Glance China Dashboard
After=network.target

[Service]
Type=simple
User=glance
Group=glance
WorkingDirectory=/opt/glance-china
ExecStart=/opt/glance-china/bin/glance-china --config /opt/glance-china/config/glance.yml
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable glance-china
sudo systemctl start glance-china
\`\`\`

### 3. Kubernetes 部署

\`\`\`yaml
# k8s-deployment.yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: glance-china
  labels:
    app: glance-china
spec:
  replicas: 2
  selector:
    matchLabels:
      app: glance-china
  template:
    metadata:
      labels:
        app: glance-china
    spec:
      containers:
      - name: glance-china
        image: glance-china:latest
        ports:
        - containerPort: 8080
        env:
        - name: REDIS_URL
          value: "redis://redis-service:6379"
        volumeMounts:
        - name: config
          mountPath: /app/config
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
      volumes:
      - name: config
        configMap:
          name: glance-config

---
apiVersion: v1
kind: Service
metadata:
  name: glance-china-service
spec:
  selector:
    app: glance-china
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
\`\`\`

## 反向代理配置

### Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
