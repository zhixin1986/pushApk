#!/bin/bash

# APK管理工具测试脚本

set -e

echo "=== APK管理工具测试 ==="

# 检查是否构建了工具
if [ ! -f "./apk-manager" ]; then
    echo "工具未构建，开始构建..."
    go build -o apk-manager .
fi

echo "1. 测试帮助信息"
./apk-manager -h

echo ""
echo "2. 测试依赖检查"
echo "检查adb是否可用："
if command -v adb &> /dev/null; then
    echo "✓ adb 可用"
    adb version | head -1
else
    echo "✗ adb 不可用"
fi

echo ""
echo "检查aapt是否可用："
if command -v aapt &> /dev/null; then
    echo "✓ aapt 可用"
    aapt version 2>/dev/null || echo "aapt 版本信息不可用"
else
    echo "✗ aapt 不可用"
fi

echo ""
echo "3. 测试设备连接"
echo "检查连接的设备："
adb devices

echo ""
echo "4. 测试完成"
echo "如果要测试完整功能，请使用："
echo "  ./apk-manager -apk /path/to/your/app.apk -v"
