# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o glance-china ./cmd/glance-china

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/glance-china .

# 创建配置目录
RUN mkdir -p /app/config /app/assets

# 复制示例配置
COPY examples/glance-china.yml /app/config/

# 设置时区为中国
ENV TZ=Asia/Shanghai

EXPOSE 8080

CMD ["./glance-china", "-config", "/app/config/glance-china.yml"]
