#!/bin/bash

# APK管理工具构建脚本

set -e

echo "开始构建APK管理工具..."

# 项目信息
PROJECT_NAME="pushApk"
VERSION="1.0.0"
BUILD_DIR="build"

# 清理构建目录
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# 构建不同平台的版本
echo "构建 Linux AMD64 版本..."
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${BUILD_DIR}/${PROJECT_NAME}-linux-amd64 .

#echo "构建 Linux ARM64 版本..."
#GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ${BUILD_DIR}/${PROJECT_NAME}-linux-arm64 .

echo "构建 Windows AMD64 版本..."
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${BUILD_DIR}/${PROJECT_NAME}-windows-amd64.exe .

echo "构建 macOS AMD64 版本..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ${BUILD_DIR}/${PROJECT_NAME}-darwin-amd64 .

#echo "构建 macOS ARM64 版本..."
#GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ${BUILD_DIR}/${PROJECT_NAME}-darwin-arm64 .

# 创建版本信息文件
echo "版本: ${VERSION}" > ${BUILD_DIR}/version.txt
echo "构建时间: $(date)" >> ${BUILD_DIR}/version.txt
echo "Git提交: $(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')" >> ${BUILD_DIR}/version.txt

echo "构建完成！"
echo "输出目录: ${BUILD_DIR}"
ls -la ${BUILD_DIR}

echo ""
echo "使用方法:"
echo "  Linux/macOS: ./${PROJECT_NAME}-linux-amd64 -apk /path/to/app.apk"
echo "  Windows: ${PROJECT_NAME}-windows-amd64.exe -apk C:\\path\\to\\app.apk"
