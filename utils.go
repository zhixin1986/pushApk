package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// checkDependencies 检查必要的工具是否存在
func checkDependencies() error {
	tools := []string{"adb", "aapt"}
	
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("找不到必要工具: %s，请确保Android SDK已正确安装并添加到PATH", tool)
		}
	}
	
	return nil
}

// checkDeviceConnection 检查设备连接
func checkDeviceConnection() error {
	cmd := exec.Command("adb", "devices")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("无法执行adb devices命令: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	deviceCount := 0
	for _, line := range lines {
		if strings.Contains(line, "device") && !strings.Contains(line, "List of devices") {
			deviceCount++
		}
	}

	if deviceCount == 0 {
		return fmt.Errorf("没有检测到连接的设备，请确保设备已连接并启用USB调试")
	}

	fmt.Printf("检测到 %d 个设备\n", deviceCount)
	return nil
}

// confirmAction 确认操作
func confirmAction(message string) bool {
	fmt.Printf("%s (y/N): ", message)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// getSystemInfo 获取系统信息
func getSystemInfo() {
	fmt.Printf("系统信息:\n")
	fmt.Printf("  操作系统: %s\n", runtime.GOOS)
	fmt.Printf("  架构: %s\n", runtime.GOARCH)
	fmt.Printf("  Go版本: %s\n", runtime.Version())
	
	// 检查adb版本
	if cmd := exec.Command("adb", "version"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			lines := strings.Split(string(output), "\n")
			if len(lines) > 0 {
				fmt.Printf("  ADB版本: %s\n", strings.TrimSpace(lines[0]))
			}
		}
	}
	
	// 检查aapt版本
	if cmd := exec.Command("aapt", "version"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			fmt.Printf("  AAPT版本: %s\n", strings.TrimSpace(string(output)))
		}
	}
	
	fmt.Println()
}

// createBackup 创建备份
func createBackup(packageName string) error {
	if !confirmAction("是否创建应用数据备份？") {
		return nil
	}

	backupDir := filepath.Join(os.TempDir(), "apk_backup")
	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		return fmt.Errorf("创建备份目录失败: %v", err)
	}

	backupFile := filepath.Join(backupDir, packageName+".ab")
	cmd := exec.Command("adb", "backup", "-apk", "-shared", "-nosystem", packageName)
	
	fmt.Printf("创建备份到: %s\n", backupFile)
	fmt.Println("请在设备上确认备份操作...")
	
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("备份失败: %v", err)
	}

	fmt.Println("备份创建成功")
	return nil
}
