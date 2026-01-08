@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

:: 构建脚本 (Windows CMD)
:: 该脚本将构建前端和后端，并将所有产物放入 build 文件夹

:: 项目根目录
set "PROJECT_ROOT=%~dp0"
set "PROJECT_ROOT=%PROJECT_ROOT:~0,-1%"
set "BUILD_DIR=%PROJECT_ROOT%\build"
set "CLIENT_DIR=%PROJECT_ROOT%\client"
set "SERVER_DIR=%PROJECT_ROOT%\server"

:: 从 build.yml 读取配置
set "BUILD_CONFIG=%PROJECT_ROOT%\build.yml"
if not exist "%BUILD_CONFIG%" (
    echo 错误: 未找到 build.yml 配置文件
    exit /b 1
)

:: 使用 PowerShell 解析 YAML 配置
for /f "delims=" %%i in ('powershell -Command "(Get-Content '%BUILD_CONFIG%' | Select-String 'image_name:').Line -replace 'image_name:\s*', ''"') do set "IMAGE_NAME=%%i"
for /f "delims=" %%i in ('powershell -Command "(Get-Content '%BUILD_CONFIG%' | Select-String 'version:').Line -replace 'version:\s*[`"'\'']*([^`"'\'']*)[`"'\'']*', '$1'"') do set "IMAGE_VERSION=%%i"
set "IMAGE_TAG=%IMAGE_NAME%:%IMAGE_VERSION%"

echo ========================================
echo   构建脚本
echo ========================================

:: 清理旧的构建目录
echo [1/6] 清理旧的构建目录...
if exist "%BUILD_DIR%" rmdir /s /q "%BUILD_DIR%"
mkdir "%BUILD_DIR%"

:: 构建前端
echo [2/6] 构建前端...
cd /d "%CLIENT_DIR%"

:: 检查是否安装了 pnpm
where pnpm >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: 未找到 pnpm，请先安装 pnpm
    exit /b 1
)

:: 安装依赖并构建
call pnpm install
if %errorlevel% neq 0 (
    echo 错误: pnpm install 失败
    exit /b 1
)

call pnpm run build:antd
if %errorlevel% neq 0 (
    echo 错误: 前端构建失败
    exit /b 1
)

:: 检查前端构建产物
set "FRONTEND_DIST=%CLIENT_DIR%\apps\web-antd\dist"
if not exist "%FRONTEND_DIST%" (
    echo 错误: 前端构建失败，未找到 dist 目录
    exit /b 1
)

:: 复制前端产物到 pb_public
echo [3/6] 复制前端产物到 pb_public...
mkdir "%BUILD_DIR%\pb_public"
xcopy /s /e /q /y "%FRONTEND_DIST%\*" "%BUILD_DIR%\pb_public\"

:: 解压 server/data.zip 到 pb_data
echo [4/6] 解压 pb_data...
mkdir "%BUILD_DIR%\pb_data"
if exist "%SERVER_DIR%\data.zip" (
    powershell -Command "Expand-Archive -Path '%SERVER_DIR%\data.zip' -DestinationPath '%BUILD_DIR%\pb_data' -Force"
) else (
    echo 警告: 未找到 server/data.zip，跳过 pb_data 解压
)

:: 复制 config 文件夹
echo [5/6] 复制 config 文件夹...
mkdir "%BUILD_DIR%\config"
xcopy /s /e /q /y "%SERVER_DIR%\config\*" "%BUILD_DIR%\config\"

:: 编译 Go 后端 (Linux amd64)
echo [6/7] 编译 Go 后端 (Linux amd64)...
cd /d "%SERVER_DIR%"

:: 检查是否安装了 Go
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: 未找到 Go，请先安装 Go
    exit /b 1
)

:: 交叉编译 Linux amd64 版本
echo 编译目标: linux/amd64
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o main .
if %errorlevel% neq 0 (
    echo 错误: Go 编译失败
    exit /b 1
)
if not exist "%SERVER_DIR%\main" (
    echo 错误: Go 编译失败，未找到 main 文件
    exit /b 1
)
echo Go 编译完成

:: 构建 Docker 镜像
echo [7/7] 构建 Docker 镜像...

:: 检查是否安装了 Docker
where docker >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: 未找到 Docker，请先安装 Docker
    exit /b 1
)

:: 构建 Docker 镜像
docker build -t %IMAGE_TAG% .
if %errorlevel% neq 0 (
    echo 错误: Docker 镜像构建失败
    exit /b 1
)

:: 保存 Docker 镜像为 tar.gz 文件
echo 保存 Docker 镜像...
docker save %IMAGE_TAG% | gzip > "%BUILD_DIR%\server-%IMAGE_VERSION%.tar.gz"
if %errorlevel% neq 0 (
    echo 错误: Docker 镜像保存失败
    exit /b 1
)

:: 复制 docker-compose.yml 并更新镜像名称
echo 复制并更新 docker-compose.yml...
powershell -Command "(Get-Content '%SERVER_DIR%\docker-compose.yml') -replace 'pocketbase-ruoyi:latest', '%IMAGE_TAG%' -replace 'pocketbase-ruoyi', 'eve-missions-pb' | Set-Content '%BUILD_DIR%\docker-compose.yml'"

:: 构建完成
cd /d "%PROJECT_ROOT%"
echo.
echo ========================================
echo   构建完成！
echo ========================================
echo 构建产物位于: %BUILD_DIR%
echo.
echo 产物列表:
echo   - config/              (配置文件)
echo   - pb_data/             (数据库文件)
echo   - pb_public/           (前端静态资源)
echo   - server-%IMAGE_VERSION%.tar.gz  (Docker 镜像)
echo   - docker-compose.yml   (Docker Compose 配置)
echo.
echo 部署说明:
echo 1. 将 build 目录中的所有文件复制到目标服务器
echo 2. 加载 Docker 镜像: docker load ^< server-%IMAGE_VERSION%.tar.gz
echo 3. 启动服务: docker-compose up -d

endlocal
