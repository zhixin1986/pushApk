# APK管理工具

一个用Go语言开发的命令行工具，用于读取APK包名、获取应用在系统中的位置、推送APK到对应位置并解压SO库。

## 功能特性

- 🔍 **读取APK包名**: 自动解析APK包名
- 📍 **获取应用位置**: 通过adb命令获取应用在系统中的安装路径
- 📤 **推送APK**: 将APK文件推送到设备的对应位置
- 📦 **解压SO库**: 自动提取APK中的SO库文件并推送到设备
- 🔒 **备份功能**: 可选择创建应用数据备份
- ✅ **依赖检查**: 自动检查必要工具是否安装
- 📱 **设备检测**: 检查设备连接状态

## 系统要求

- Go 1.21 或更高版本
- Android SDK (包含adb和aapt工具)
- 已连接的Android设备（启用USB调试）
- Root权限（用于系统级操作）

## 安装

### 1. 克隆项目
```bash
git clone <repository-url>
cd pushApk
```

### 2. 构建
```bash
go build -o pushApk .
```

### 3. 安装（可选）
```bash
# Linux/macOS
sudo mv pushApk /usr/local/bin/

# Windows
# 将pushApk.exe移动到PATH中的某个目录
```

## 使用方法

### 基本使用
```bash
pushApk /path/to/your/app.apk
pushApk -apk /path/to/your/app.apk
```

### 高级选项
```bash
# 详细输出
pushApk -apk /path/to/your/app.apk -v

# 创建备份
pushApk -apk /path/to/your/app.apk -backup

# 组合使用
pushApk -apk /path/to/your/app.apk -backup -v
```

### 命令行选项

| 选项 | 描述 |
|------|------|
| `-apk` | APK文件路径（必需） |
| `-v` | 详细输出模式 |
| `-backup` | 创建应用数据备份 |
| `-h` | 显示帮助信息 |

## 工作流程

1. **依赖检查**: 确保adb和aapt工具可用
2. **设备检测**: 检查设备连接状态
3. **APK解析**: 读取APK文件并提取包名
4. **路径获取**: 获取应用在系统中的安装路径
5. **备份创建**: （可选）创建应用数据备份
6. **APK推送**: 将APK文件推送到设备
7. **SO库处理**: 提取并推送SO库文件

## 前置条件

### Android SDK 安装
确保已安装Android SDK并将tools和platform-tools目录添加到PATH：

```bash
# 添加到 ~/.bashrc 或 ~/.zshrc
export ANDROID_HOME=/path/to/android-sdk
export PATH=$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/tools
```

### 设备准备
1. 启用开发者选项
2. 开启USB调试
3. 获取Root权限（用于系统级操作）
4. 连接设备并授权调试

### 验证环境
```bash
# 检查adb
adb version

# 检查aapt
aapt version

# 检查设备连接
adb devices
```

## 注意事项

⚠️ **重要警告**:
- 此工具需要Root权限才能正常工作
- 操作系统应用可能会影响设备稳定性
- 建议在操作前创建完整的设备备份
- 仅在测试设备上使用，避免在生产设备上使用

## 故障排除

### 常见问题

1. **找不到adb/aapt工具**
   - 确保Android SDK已正确安装
   - 检查PATH环境变量

2. **设备未连接**
   - 检查USB连接
   - 确保USB调试已启用
   - 重新授权调试权限

3. **权限不足**
   - 确保设备已获取Root权限
   - 检查adb是否以root身份运行

4. **APK解析失败**
   - 确保APK文件完整且未损坏
   - 检查aapt工具版本

## 示例输出

```
系统信息:
  操作系统: linux
  架构: amd64
  Go版本: go1.21.0
  ADB版本: Android Debug Bridge version 1.0.41
  AAPT版本: Android Asset Packaging Tool

检测到 1 个设备
开始处理APK文件: /path/to/app.apk
包名: com.example.app
应用路径: /data/app/com.example.app
确定要继续推送APK和SO库吗？ (y/N): y
APK推送成功
推送SO文件成功: lib/arm64-v8a/libnative.so
推送SO文件成功: lib/armeabi-v7a/libnative.so
SO库解压成功
处理完成!
```

## 许可证

MIT License

## 贡献

欢迎提交Issues和Pull Requests！

## 更新日志

### v1.0.0
- 初始版本
- 支持APK包名读取
- 支持应用路径获取
- 支持APK推送
- 支持SO库解压和推送
- 添加备份功能
- 添加依赖检查
