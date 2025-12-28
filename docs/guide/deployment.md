# 部署指南

## 部署方式

PocketBase RuoYi 支持多种部署方式，你可以根据需求选择合适的方案。

## Docker 部署

### 使用 Docker Compose

项目已包含 `docker-compose.yml` 文件，可以快速部署：

```bash
cd server
docker-compose up -d
```

### 自定义 Dockerfile

如需自定义构建，可以修改 `server/Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/pb_data ./pb_data

EXPOSE 8090
CMD ["./main", "serve", "--http=0.0.0.0:8090"]
```

构建和运行：

```bash
docker build -t pocketbase-ruoyi .
docker run -d -p 8090:8090 -v $(pwd)/pb_data:/app/pb_data pocketbase-ruoyi
```

## 传统部署

### 后端部署

1. 构建 Go 程序：

```bash
cd server
go build -o main
```

2. 上传到服务器并运行：

```bash
# 直接运行
./main serve --http=0.0.0.0:8090

# 或使用 systemd 服务
sudo nano /etc/systemd/system/pocketbase.service
```

systemd 服务配置示例：

```ini
[Unit]
Description=PocketBase Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/pocketbase
ExecStart=/var/www/pocketbase/main serve --http=0.0.0.0:8090
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable pocketbase
sudo systemctl start pocketbase
```

### 前端部署

1. 构建前端资源：

```bash
cd client
pnpm build
```

2. 部署到 Web 服务器（Nginx 配置示例）：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /var/www/pocketbase-ruoyi/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 云平台部署

### Vercel 部署（前端）

1. 安装 Vercel CLI：

```bash
npm i -g vercel
```

2. 部署：

```bash
cd client
vercel --prod
```

### Railway / Render（后端）

这些平台支持直接从 Git 仓库部署 Go 应用：

1. 连接 GitHub 仓库
2. 选择 `server` 目录作为根目录
3. 设置构建命令：`go build -o main`
4. 设置启动命令：`./main serve --http=0.0.0.0:$PORT`

## 环境变量配置

确保在生产环境设置正确的环境变量：

```env
# 生产环境配置
APP_ENV=production
APP_URL=https://your-domain.com

# 数据库备份
BACKUP_ENABLED=true
BACKUP_CRON=0 2 * * *

# 邮件配置
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your-email@example.com
SMTP_PASSWORD=your-password
```

## 数据备份

### 自动备份

PocketBase 支持自动备份：

```bash
./main serve --backup-cron="0 2 * * *"
```

### 手动备份

```bash
# 备份数据库
cp -r pb_data pb_data_backup_$(date +%Y%m%d)

# 或使用 PocketBase 内置备份
./main backup
```

## 性能优化

### 启用 Gzip

在 Nginx 配置中启用 Gzip 压缩：

```nginx
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript application/javascript application/json;
```

### 启用缓存

```nginx
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## 安全建议

1. **使用 HTTPS**：配置 SSL 证书（推荐 Let's Encrypt）
2. **设置强密码**：Admin 账户使用强密码
3. **限制 Admin UI 访问**：通过防火墙或 Nginx 限制 `/_/` 路径访问
4. **定期备份**：设置自动备份策略
5. **更新依赖**：定期更新 PocketBase 和前端依赖

## 监控与日志

### 日志查看

```bash
# systemd 服务日志
sudo journalctl -u pocketbase -f

# Docker 日志
docker logs -f pocketbase-ruoyi
```

### 监控工具

可以集成以下监控工具：

- Prometheus + Grafana
- Sentry（错误追踪）
- Uptime Kuma（服务监控）
