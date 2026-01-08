# 构建脚本 (PowerShell)
# 该脚本将构建前端和后端，并将所有产物放入 build 文件夹

$ErrorActionPreference = "Stop"

# 项目根目录
$PROJECT_ROOT = Split-Path -Parent $MyInvocation.MyCommand.Path
$BUILD_DIR = Join-Path $PROJECT_ROOT "build"
$CLIENT_DIR = Join-Path $PROJECT_ROOT "client"
$SERVER_DIR = Join-Path $PROJECT_ROOT "server"

# 从 build.yml 读取配置
$BUILD_CONFIG = Join-Path $PROJECT_ROOT "build.yml"
if (-not (Test-Path $BUILD_CONFIG)) {
    Write-Host "错误: 未找到 build.yml 配置文件" -ForegroundColor Red
    exit 1
}

$ConfigContent = Get-Content $BUILD_CONFIG -Raw
if ($ConfigContent -match "image_name:\s*(.+)") {
    $IMAGE_NAME = $matches[1].Trim()
}
if ($ConfigContent -match "version:\s*[`"']?([^`"'\r\n]+)[`"']?") {
    $IMAGE_VERSION = $matches[1].Trim()
}
$IMAGE_TAG = "${IMAGE_NAME}:${IMAGE_VERSION}"

Write-Host "========================================" -ForegroundColor Green
Write-Host "  构建脚本" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# 清理旧的构建目录
Write-Host "[1/6] 清理旧的构建目录..." -ForegroundColor Yellow
if (Test-Path $BUILD_DIR) {
    Remove-Item -Recurse -Force $BUILD_DIR
}
New-Item -ItemType Directory -Path $BUILD_DIR -Force | Out-Null

# 构建前端
Write-Host "[2/6] 构建前端..." -ForegroundColor Yellow
Set-Location $CLIENT_DIR

# 检查是否安装了 pnpm
if (-not (Get-Command pnpm -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到 pnpm，请先安装 pnpm" -ForegroundColor Red
    exit 1
}

# 安装依赖并构建
pnpm install
pnpm run build:antd

# 检查前端构建产物
$FRONTEND_DIST = Join-Path $CLIENT_DIR "apps\web-antd\dist"
if (-not (Test-Path $FRONTEND_DIST)) {
    Write-Host "错误: 前端构建失败，未找到 dist 目录" -ForegroundColor Red
    exit 1
}

# 复制前端产物到 pb_public
Write-Host "[3/6] 复制前端产物到 pb_public..." -ForegroundColor Yellow
$PB_PUBLIC = Join-Path $BUILD_DIR "pb_public"
New-Item -ItemType Directory -Path $PB_PUBLIC -Force | Out-Null
Copy-Item -Recurse -Force "$FRONTEND_DIST\*" $PB_PUBLIC

# 解压 server/data.zip 到 pb_data
Write-Host "[4/6] 解压 pb_data..." -ForegroundColor Yellow
$PB_DATA = Join-Path $BUILD_DIR "pb_data"
New-Item -ItemType Directory -Path $PB_DATA -Force | Out-Null
$DATA_ZIP = Join-Path $SERVER_DIR "data.zip"
if (Test-Path $DATA_ZIP) {
    Expand-Archive -Path $DATA_ZIP -DestinationPath $PB_DATA -Force
} else {
    Write-Host "警告: 未找到 server/data.zip，跳过 pb_data 解压" -ForegroundColor Yellow
}

# 复制 config 文件夹
Write-Host "[5/6] 复制 config 文件夹..." -ForegroundColor Yellow
$CONFIG_DIR = Join-Path $BUILD_DIR "config"
New-Item -ItemType Directory -Path $CONFIG_DIR -Force | Out-Null
Copy-Item -Recurse -Force (Join-Path $SERVER_DIR "config\*") $CONFIG_DIR

# 编译 Go 后端 (Linux amd64)
Write-Host "[6/7] 编译 Go 后端 (Linux amd64)..." -ForegroundColor Yellow
Set-Location $SERVER_DIR

# 检查是否安装了 Go
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到 Go，请先安装 Go" -ForegroundColor Red
    exit 1
}

# 交叉编译 Linux amd64 版本
Write-Host "编译目标: linux/amd64"
$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o main .
if (-not (Test-Path (Join-Path $SERVER_DIR "main"))) {
    Write-Host "错误: Go 编译失败" -ForegroundColor Red
    exit 1
}
Write-Host "Go 编译完成" -ForegroundColor Green

# 构建 Docker 镜像
Write-Host "[7/7] 构建 Docker 镜像..." -ForegroundColor Yellow

# 检查是否安装了 Docker
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到 Docker，请先安装 Docker" -ForegroundColor Red
    exit 1
}

# 构建 Docker 镜像
docker build -t $IMAGE_TAG .

# 保存 Docker 镜像为 tar.gz 文件
Write-Host "保存 Docker 镜像..." -ForegroundColor Yellow
$TAR_FILE = Join-Path $BUILD_DIR "server-${IMAGE_VERSION}.tar.gz"
docker save $IMAGE_TAG | gzip > $TAR_FILE

# 复制 docker-compose.yml 并更新镜像名称
Write-Host "复制并更新 docker-compose.yml..." -ForegroundColor Yellow
$COMPOSE_CONTENT = Get-Content (Join-Path $SERVER_DIR "docker-compose.yml") -Raw
$COMPOSE_CONTENT = $COMPOSE_CONTENT -replace "pocketbase-ruoyi:latest", $IMAGE_TAG
$COMPOSE_CONTENT = $COMPOSE_CONTENT -replace "pocketbase-ruoyi", "eve-missions-pb"
$COMPOSE_CONTENT | Set-Content (Join-Path $BUILD_DIR "docker-compose.yml") -NoNewline

# 构建完成
Set-Location $PROJECT_ROOT
Write-Host "========================================" -ForegroundColor Green
Write-Host "  构建完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host "构建产物位于: $BUILD_DIR" -ForegroundColor Green
Write-Host ""
Write-Host "产物列表:"
Write-Host "  - config/              (配置文件)"
Write-Host "  - pb_data/             (数据库文件)"
Write-Host "  - pb_public/           (前端静态资源)"
Write-Host "  - server-${IMAGE_VERSION}.tar.gz  (Docker 镜像)"
Write-Host "  - docker-compose.yml   (Docker Compose 配置)"
Write-Host ""
Write-Host "部署说明:" -ForegroundColor Yellow
Write-Host "1. 将 build 目录中的所有文件复制到目标服务器"
Write-Host "2. 加载 Docker 镜像: docker load < server-${IMAGE_VERSION}.tar.gz"
Write-Host "3. 启动服务: docker-compose up -d"
