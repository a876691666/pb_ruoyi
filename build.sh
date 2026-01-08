#!/bin/bash

# 构建脚本
# 该脚本将构建前端和后端，并将所有产物放入 build 文件夹

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
BUILD_DIR="$PROJECT_ROOT/build"
CLIENT_DIR="$PROJECT_ROOT/client"
SERVER_DIR="$PROJECT_ROOT/server"

# 从 build.yml 读取配置
BUILD_CONFIG="$PROJECT_ROOT/build.yml"
if [ ! -f "$BUILD_CONFIG" ]; then
    echo -e "${RED}错误: 未找到 build.yml 配置文件${NC}"
    exit 1
fi

# 解析 YAML 配置（简单解析）
IMAGE_NAME=$(grep 'image_name:' "$BUILD_CONFIG" | sed 's/image_name: *//g' | tr -d '\r')
IMAGE_VERSION=$(grep 'version:' "$BUILD_CONFIG" | sed 's/version: *//g' | tr -d '"' | tr -d "'" | tr -d '\r')
IMAGE_TAG="$IMAGE_NAME:$IMAGE_VERSION"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  构建脚本${NC}"
echo -e "${GREEN}========================================${NC}"

# 清理旧的构建目录
echo -e "${YELLOW}[1/6] 清理旧的构建目录...${NC}"
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# 构建前端
echo -e "${YELLOW}[2/6] 构建前端...${NC}"
cd "$CLIENT_DIR"

# 检查是否安装了 pnpm
if ! command -v pnpm &> /dev/null; then
    echo -e "${RED}错误: 未找到 pnpm，请先安装 pnpm${NC}"
    exit 1
fi

# 安装依赖并构建
pnpm install
pnpm run build:antd

# 检查前端构建产物
FRONTEND_DIST="$CLIENT_DIR/apps/web-antd/dist"
if [ ! -d "$FRONTEND_DIST" ]; then
    echo -e "${RED}错误: 前端构建失败，未找到 dist 目录${NC}"
    exit 1
fi

# 复制前端产物到 pb_public
echo -e "${YELLOW}[3/6] 复制前端产物到 pb_public...${NC}"
mkdir -p "$BUILD_DIR/pb_public"
cp -r "$FRONTEND_DIST"/* "$BUILD_DIR/pb_public/"

# 解压 server/data.zip 到 pb_data
echo -e "${YELLOW}[4/6] 解压 pb_data...${NC}"
mkdir -p "$BUILD_DIR/pb_data"
if [ -f "$SERVER_DIR/data.zip" ]; then
    unzip -o "$SERVER_DIR/data.zip" -d "$BUILD_DIR/pb_data/"
else
    echo -e "${RED}警告: 未找到 server/data.zip，跳过 pb_data 解压${NC}"
fi

# 复制 config 文件夹
echo -e "${YELLOW}[5/6] 复制 config 文件夹...${NC}"
mkdir -p "$BUILD_DIR/config"
cp -r "$SERVER_DIR/config"/* "$BUILD_DIR/config/"

# 编译 Go 后端 (Linux amd64)
echo -e "${YELLOW}[6/7] 编译 Go 后端 (Linux amd64)...${NC}"
cd "$SERVER_DIR"

# 检查是否安装了 Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: 未找到 Go，请先安装 Go${NC}"
    exit 1
fi

# 交叉编译 Linux amd64 版本
echo -e "编译目标: linux/amd64"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main .
if [ ! -f "$SERVER_DIR/main" ]; then
    echo -e "${RED}错误: Go 编译失败${NC}"
    exit 1
fi
echo -e "${GREEN}Go 编译完成${NC}"

# 构建 Docker 镜像
echo -e "${YELLOW}[7/7] 构建 Docker 镜像...${NC}"

# 检查是否安装了 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: 未找到 Docker，请先安装 Docker${NC}"
    exit 1
fi

# 构建 Docker 镜像
docker build -t "$IMAGE_TAG" .

# 保存 Docker 镜像为 tar.gz 文件
echo -e "${YELLOW}保存 Docker 镜像...${NC}"
docker save "$IMAGE_TAG" | gzip > "$BUILD_DIR/server-$IMAGE_VERSION.tar.gz"

# 复制 docker-compose.yml 并更新镜像名称
echo -e "${YELLOW}复制并更新 docker-compose.yml...${NC}"
sed "s|pocketbase-ruoyi:latest|$IMAGE_TAG|g" "$SERVER_DIR/docker-compose.yml" | \
sed "s|pocketbase-ruoyi|eve-missions-pb|g" > "$BUILD_DIR/docker-compose.yml"

# 构建完成
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  构建完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}构建产物位于: $BUILD_DIR${NC}"
echo ""
echo -e "产物列表:"
echo -e "  - config/              (配置文件)"
echo -e "  - pb_data/             (数据库文件)"
echo -e "  - pb_public/           (前端静态资源)"
echo -e "  - server-$IMAGE_VERSION.tar.gz  (Docker 镜像)"
echo -e "  - docker-compose.yml   (Docker Compose 配置)"
echo ""
echo -e "${YELLOW}部署说明:${NC}"
echo -e "1. 将 build 目录中的所有文件复制到目标服务器"
echo -e "2. 加载 Docker 镜像: docker load < server-$IMAGE_VERSION.tar.gz"
echo -e "3. 启动服务: docker-compose up -d"
