package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/shogo82148/androidbinary/apk"
)

type APKManager struct {
	adbPath string
}

func NewAPKManager() *APKManager {
	return &APKManager{
		adbPath: "adb", // 假设adb在PATH中
	}
}

// GetPackageName 从APK文件中提取包名
func (m *APKManager) GetPackageName(apkPath string) (string, error) {
	// 使用androidbinary库解析APK文件
	apkFile, err := apk.OpenFile(apkPath)
	if err != nil {
		return "", fmt.Errorf("打开APK文件失败: %v", err)
	}
	defer apkFile.Close()

	// 获取manifest信息
	manifest := apkFile.Manifest()

	// 提取包名
	packageName := manifest.Package.MustString()
	if packageName == "" {
		return "", fmt.Errorf("无法从APK中提取包名")
	}

	return packageName, nil
}

// GetAppPath 获取应用在系统中的安装路径
func (m *APKManager) GetAppPath(packageName string) (string, error) {
	// 使用adb命令获取应用路径
	cmd := exec.Command(m.adbPath, "shell", "pm", "path", packageName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取应用路径失败: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		return "", fmt.Errorf("应用未安装: %s", packageName)
	}

	// 解析输出，格式通常是 "package:/data/app/package.name/base.apk"
	if strings.HasPrefix(outputStr, "package:") {
		path := strings.TrimPrefix(outputStr, "package:")
		return path, nil
	}

	return "", fmt.Errorf("无法解析应用路径: %s", outputStr)
}

// PushAPK 推送APK到指定位置
func (m *APKManager) PushAPK(apkPath, targetPath string) error {
	// 获取APK文件名
	targetFile := targetPath
	// 使用adb push命令
	cmd := exec.Command(m.adbPath, "push", apkPath, targetFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("推送APK失败: %v, 输出: %s", err, string(output))
	}

	fmt.Printf("推送输出: %s\n", string(output))
	return nil
}

func (m *APKManager) Reboot() error {
	cmd := exec.Command(m.adbPath, "shell", "sync")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("同步数据失败: %v, 输出: %s", err, string(output))
	}
	cmd = exec.Command(m.adbPath, "shell", "stop")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("停止系统服务失败: %v, 输出: %s", err, string(output))
	}
	cmd = exec.Command(m.adbPath, "shell", "start")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("启动系统服务失败: %v, 输出: %s", err, string(output))
	}
	return nil
}

// ExtractSOLibraries 解压APK中的SO库
func (m *APKManager) ExtractSOLibraries(apkPath, targetPath string) error {
	// 打开APK文件（APK本质上是ZIP文件）
	reader, err := zip.OpenReader(apkPath)
	if err != nil {
		return fmt.Errorf("打开APK文件失败: %v", err)
	}
	defer reader.Close()

	// 创建临时目录用于存放SO库
	tempDir := filepath.Join(os.TempDir(), "apk_so_libs")
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir) // 确保清理

	// 分析APK中的SO库
	soFiles := make(map[string][]string) // arch -> []soFiles
	for _, file := range reader.File {
		if strings.HasPrefix(file.Name, "lib/") && strings.HasSuffix(file.Name, ".so") {
			// 提取架构信息，例如 lib/arm64-v8a/libnative.so
			parts := strings.Split(file.Name, "/")
			if len(parts) >= 3 {
				arch := parts[1]
				if soFiles[arch] == nil {
					soFiles[arch] = make([]string, 0)
				}
				soFiles[arch] = append(soFiles[arch], file.Name)
			}
		}
	}

	if len(soFiles) == 0 {
		fmt.Println("APK中没有找到SO库文件")
		return nil
	}

	// 显示找到的SO库信息
	fmt.Printf("发现SO库架构: %v\n", getKeys(soFiles))

	// 获取设备架构
	deviceArch, err := m.getDeviceArchitecture()
	if err != nil {
		fmt.Printf("获取设备架构失败: %v, 将尝试推送所有架构的SO库\n", err)
		deviceArch = ""
	} else {
		fmt.Printf("设备架构: %s\n", deviceArch)
	}

	libDir := filepath.Join(targetPath, "lib")
	fmt.Printf("删除目录 %s\n", libDir)
	cmd := exec.Command(m.adbPath, "shell", "rm -rf "+libDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("删除目录失败: %v, 输出: %s", err, string(output))
	}
	// 提取并推送SO库
	for arch, files := range soFiles {
		// 如果知道设备架构，优先推送对应架构的SO库
		if deviceArch != "" && arch != deviceArch {
			continue
		}

		fmt.Printf("处理架构 %s 的SO库 (%d个文件)\n", arch, len(files))
		for _, soFile := range files {
			// 提取SO文件到临时目录
			err := m.extractFileFromZip(reader, soFile, tempDir)
			if err != nil {
				fmt.Printf("提取SO文件失败 %s: %v\n", soFile, err)
				continue
			}

			// 推送SO文件到设备
			localPath := filepath.Join(tempDir, soFile)
			remotePath := filepath.Join(targetPath, "lib", arch, filepath.Base(soFile))

			// 确保远程目录存在
			remoteDir := filepath.Dir(remotePath)

			//转换成linux路径
			remotePath = filepath.ToSlash(remotePath)
			remoteDir = filepath.ToSlash(remoteDir)
			cmd = exec.Command(m.adbPath, "shell", "mkdir -p "+remoteDir)
			cmd.Run()

			// 推送SO文件
			cmd = exec.Command(m.adbPath, "push", localPath, remotePath)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("推送SO文件失败 %s: %v, 输出: %s\n", soFile, err, string(output))
				continue
			}

			fmt.Printf("推送SO文件成功: %s -> %s\n", soFile, remotePath)
		}

		// 如果找到了设备对应的架构，就不处理其他架构了
		if deviceArch != "" && arch == deviceArch {
			break
		}
	}

	return nil
}

// getDeviceArchitecture 获取设备架构
func (m *APKManager) getDeviceArchitecture() (string, error) {
	cmd := exec.Command(m.adbPath, "shell", "getprop", "ro.product.cpu.abi")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	arch := strings.TrimSpace(string(output))
	// 映射到标准架构名称
	switch arch {
	case "arm64-v8a":
		return "arm64-v8a", nil
	case "armeabi-v7a":
		return "armeabi-v7a", nil
	case "x86_64":
		return "x86_64", nil
	case "x86":
		return "x86", nil
	default:
		return arch, nil
	}
}

// extractFileFromZip 从ZIP中提取指定文件
func (m *APKManager) extractFileFromZip(reader *zip.ReadCloser, fileName, destDir string) error {
	for _, file := range reader.File {
		if file.Name == fileName {
			return m.extractFile(file, destDir)
		}
	}
	return fmt.Errorf("文件未找到: %s", fileName)
}

// getKeys 获取map的所有键
func getKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// extractFile 从ZIP中提取单个文件
func (m *APKManager) extractFile(file *zip.File, destDir string) error {
	// 创建目标文件路径
	destPath := filepath.Join(destDir, file.Name)

	// 确保目标目录存在
	err := os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return err
	}

	// 打开ZIP中的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	// 复制文件内容
	_, err = io.Copy(dest, src)
	return err
}
