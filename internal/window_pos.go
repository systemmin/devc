/**
 * @Time : 2025/3/16 10:04
 * @File : window_pos.go
 * @Software: dev_clip
 * @Author : Mr.Fang
 * @Description: cmd 窗口置顶
 */

package internal

import (
	"log"
	"syscall"
	"unsafe"
)

// 定义 Windows API 常量
const (
	HWND_TOPMOST = ^uintptr(0) // 正确的 HWND_TOPMOST 值 (0xFFFFFFFFFFFFFFFF)
	SWP_NOSIZE   = 0x0001
	SWP_NOMOVE   = 0x0002
)

var (
	findWindow   = user32.NewProc("FindWindowW")
	setWindowPos = user32.NewProc("SetWindowPos")
)

// 获取控制台窗口句柄（适配 Windows 11）
func getConsoleWindow() syscall.Handle {
	// 尝试不同控制台类名（传统cmd和Windows Terminal）
	classNames := []string{
		"ConsoleWindowClass",            // 传统cmd/PowerShell
		"CASCADIA_HOSTING_WINDOW_CLASS", // Windows Terminal
	}

	for _, className := range classNames {
		ptrClassName := syscall.StringToUTF16Ptr(className)
		ret, _, _ := findWindow.Call(
			uintptr(unsafe.Pointer(ptrClassName)),
			uintptr(0), // 不指定窗口标题
		)
		if ret != 0 {
			log.Printf("成功找到窗口类名: %s, 句柄: 0x%x", className, ret)
			return syscall.Handle(ret)
		}
	}
	return 0
}

// TopUp 置顶
func TopUp() {
	// 获取控制台窗口句柄
	hwnd := getConsoleWindow()
	if hwnd == 0 {
		log.Fatal("未找到控制台窗口句柄，请确认是否在终端中运行程序")
	}

	// 设置窗口置顶
	ret, _, _ := setWindowPos.Call(
		uintptr(hwnd),
		HWND_TOPMOST,
		0, 0, 0, 0,
		SWP_NOSIZE|SWP_NOMOVE,
	)
	if ret == 0 {
		log.Fatal("设置窗口置顶失败，错误码:", syscall.GetLastError())
	}

	log.Println("窗口已置顶！按 Ctrl+C 退出程序")
}
