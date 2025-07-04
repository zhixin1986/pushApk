package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var apkPath string
	var help bool
	var verbose bool
	var backup bool
	var debug bool
	var dryRun bool
	var targetPath string
	var skipSO bool

	flag.StringVar(&apkPath, "apk", "", "APK文件路径")
	flag.BoolVar(&help, "h", false, "显示帮助信息")
	flag.BoolVar(&verbose, "v", false, "详细输出")
	flag.BoolVar(&backup, "backup", false, "创建备份")
	flag.BoolVar(&debug, "debug", false, "启用调试模式")
	flag.BoolVar(&dryRun, "dry-run", false, "预览模式，不实际执行操作")
	flag.StringVar(&targetPath, "target", "", "自定义目标路径")
	flag.BoolVar(&skipSO, "skip-so", false, "跳过SO库文件处理")
	flag.Parse()
	if flag.NArg() > 0 && apkPath == "" {
		apkPath = flag.Arg(0)
	}
	if help || apkPath == "" {
		showHelp()
		return
	}

	// 设置调试和详细输出
	if debug {
		verbose = true
		fmt.Println("调试模式已启用")
	}

	if verbose {
		getSystemInfo()
	}

	if dryRun {
		fmt.Println("预览模式已启用，将不执行实际操作")
	}

	// 检查设备连接
	if err := checkDeviceConnection(); err != nil {
		log.Fatalf("设备连接检查失败: %v", err)
	}

	// 检查APK文件是否存在
	if _, err := os.Stat(apkPath); os.IsNotExist(err) {
		log.Fatalf("APK文件不存在: %s", apkPath)
	}

	fmt.Printf("开始处理APK文件: %s\n", apkPath)

	// 创建APK管理器
	manager := NewAPKManager()

	// 读取APK包名
	// 获取包名
	if debug {
		fmt.Printf("正在获取APK包名: %s\n", apkPath)
	}
	packageName, err := manager.GetPackageName(apkPath)
	if err != nil {
		log.Fatalf("获取包名失败: %v", err)
	}
	fmt.Printf("包名: %s\n", packageName)

	// 创建备份（如果需要）
	if backup {
		if debug {
			fmt.Printf("正在创建备份: %s\n", packageName)
		}
		if !dryRun {
			if err := createBackup(packageName); err != nil {
				log.Printf("备份失败: %v", err)
			}
		} else {
			fmt.Printf("预览: 将创建备份 %s\n", packageName)
		}
	}

	// 获取应用在系统中的位置
	var appDir string
	var appPath string
	if targetPath != "" {
		appDir = targetPath
		if debug {
			fmt.Printf("使用自定义目标路径: %s\n", appDir)
		}
	} else {
		if debug {
			fmt.Printf("正在获取应用安装路径: %s\n", packageName)
		}
		var err error
		appPath, err = manager.GetAppPath(packageName)
		appDir = filepath.Dir(appPath)
		appDir = filepath.ToSlash(appDir)
		if err != nil {
			log.Fatalf("获取应用路径失败: %v", err)
		}
	}
	fmt.Printf("目标路径: %s\n", appDir)

	// 确认操作（预览模式跳过）
	if !dryRun && !confirmAction("确定要继续推送APK和SO库吗？") {
		fmt.Println("操作已取消")
		return
	}

	// 推送APK到对应位置

	fmt.Printf("正在推送APK: %s -> %s\n", apkPath, appPath)
	if !dryRun {
		err = manager.PushAPK(apkPath, appPath)
		if err != nil {
			log.Fatalf("推送APK失败: %v", err)
		}
		fmt.Printf("APK推送成功\n")
	} else {
		fmt.Printf("预览: 将推送APK %s 到 %s\n", apkPath, appDir)
	}

	// 解压SO库
	if !skipSO {
		if debug {
			fmt.Printf("正在处理SO库: %s -> %s\n", apkPath, appDir)
		}
		if !dryRun {
			err = manager.ExtractSOLibraries(apkPath, appDir)
			if err != nil {
				log.Fatalf("解压SO库失败: %v", err)
			}
			fmt.Printf("SO库处理成功\n")
		} else {
			fmt.Printf("预览: 将解压SO库从 %s 到 %s\n", apkPath, appDir)
		}
	} else {
		if debug {
			fmt.Println("跳过SO库处理")
		}
	}

	if dryRun {
		fmt.Println("预览模式完成!")
	} else {
		fmt.Println("处理完成!")
	}
	if confirmAction("是否需要重启设备？") {
		if !dryRun {
			err := manager.Reboot()
			if err != nil {
				log.Fatalf("重启设备失败: %v", err)
			}
		} else {
			fmt.Println("预览: 将重启设备")
		}
	}
}

func showHelp() {
	fmt.Println("APK管理工具")
	fmt.Println("用法: pushApk [-apk] <APK文件路径> [选项]")
	fmt.Println("选项:")
	fmt.Println("  -apk string      APK文件路径")
	fmt.Println("  -v               详细输出")
	fmt.Println("  -backup          创建备份")
	fmt.Println("  -debug           启用调试模式 (包含详细输出)")
	fmt.Println("  -dry-run         预览模式，不实际执行操作")
	fmt.Println("  -target string   自定义目标路径")
	fmt.Println("  -skip-so         跳过SO库文件处理")
	fmt.Println("  -h               显示帮助信息")
	fmt.Println("")
	fmt.Println("功能:")
	fmt.Println("  1. 读取APK包名")
	fmt.Println("  2. 获取应用在系统中的位置")
	fmt.Println("  3. 推送APK到对应位置")
	fmt.Println("  4. 解压并推送SO库")
	fmt.Println("")
	fmt.Println("前置条件:")
	fmt.Println("  - 安装Android SDK (adb工具)")
	fmt.Println("  - 设备已连接并启用USB调试")
	fmt.Println("  - 设备已获取Root权限")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  pushApk  /path/to/app.apk")
	fmt.Println("  pushApk -apk /path/to/app.apk")
	fmt.Println("  pushApk -apk /path/to/app.apk -backup -v")
	fmt.Println("  pushApk -apk /path/to/app.apk -debug -dry-run")
	fmt.Println("  pushApk -apk /path/to/app.apk -target /custom/path")
}
