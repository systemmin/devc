/**
 * @Time : 2025/3/15 17:22
 * @File : clipboard.go
 * @Software: dev_clip
 * @Author : Mr.Fang
 * @Description: 监听 window 剪切板
 */

package internal

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

// Windows API https://learn.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-openclipboard
var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32")

	procOpenClipboard    = user32.NewProc("OpenClipboard")    // 打开剪切板
	procCloseClipboard   = user32.NewProc("CloseClipboard")   // 关闭剪切板
	procGetClipboardData = user32.NewProc("GetClipboardData") // 获取剪切板数据
	procSetClipboardData = user32.NewProc("SetClipboardData") // 设置剪切板数据
	procEmptyClipboard   = user32.NewProc("EmptyClipboard")   // 打开剪切板
	procGlobalAlloc      = kernel32.NewProc("GlobalAlloc")    // 全局内存
	procGlobalLock       = kernel32.NewProc("GlobalLock")     // 全局锁
	procGlobalUnlock     = kernel32.NewProc("GlobalUnlock")   // 释放锁
	procRtlMoveMemory    = kernel32.NewProc("RtlMoveMemory")  // 用于拷贝数据

	getClipboardSequenceNumber     = user32.NewProc("GetClipboardSequenceNumber") // 剪切版序号
	procIsClipboardFormatAvailable = user32.NewProc("IsClipboardFormatAvailable") // 检查指定格式文本
)

const (
	GMEM_MOVEABLE  = 0x0002 // 分配可移动内存 https://learn.microsoft.com/zh-cn/windows/win32/api/winbase/nf-winbase-globalalloc#gmem_moveable
	CF_UNICODETEXT = 13     // Unicode 文本格式 https://learn.microsoft.com/zh-cn/windows/win32/dataxchg/standard-clipboard-formats#cf_unicodetext
)

// ReadClipboard 读取剪切板内容
func ReadClipboard() string {
	var content string
	// 尝试打开剪切板
	if OpenClipboard() {
		// 检查剪切板是否包含Unicode文本格式
		if IsClipboardFormatAvailable(CF_UNICODETEXT) {
			// 读取剪切板中的文本内容
			data, err := GetClipboardData(CF_UNICODETEXT)
			if err != nil {
				fmt.Println("获取剪贴板数据时出错:", err)
			} else {
				content = data
			}
		}
		// 关闭剪切板
		CloseClipboard()
	} else {
		fmt.Println("打开剪切板失败")
	}
	return content
}

// Watch 监听剪切板
func Watch() <-chan string {
	receive := make(chan string, 1)
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop() // 确保在程序退出时停止 Ticker
		cnt, _, _ := getClipboardSequenceNumber.Call()
		// 使用 select 从 ticker 的通道接收时间戳
		for {
			select {
			case t := <-ticker.C: // 接收时间戳
				cur, _, _ := getClipboardSequenceNumber.Call()
				if cnt != cur {
					format := t.Format(time.DateTime)
					clipboard := ReadClipboard()
					fmt.Println(format)
					receive <- clipboard
					cnt = cur
				}
			}
		}
	}()
	return receive
}

// OpenClipboard 打开剪切板
func OpenClipboard() bool {
	ret, _, _ := procOpenClipboard.Call()
	return ret != 0
}

// CloseClipboard 关闭剪切板
func CloseClipboard() {
	procCloseClipboard.Call()
}

// IsClipboardFormatAvailable 检查剪切板是否包含指定格式
func IsClipboardFormatAvailable(format uint) bool {
	ret, _, _ := procIsClipboardFormatAvailable.Call(uintptr(format))
	return ret != 0
}

// GetClipboardData 读取剪切板中的数据
func GetClipboardData(format uint) (string, error) {
	hMem, _, _ := procGetClipboardData.Call(uintptr(format))
	if hMem == 0 {
		return "", fmt.Errorf("未能获取剪贴板数据句柄")
	}
	// 锁定全局内存
	ptr, _, _ := procGlobalLock.Call(hMem)
	if ptr == 0 {
		return "", fmt.Errorf("无法锁定全局内存")
	}
	defer procGlobalUnlock.Call(hMem) // 释放锁定的内存

	// 读取 UTF-16 数据
	data := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])

	return data, nil
}

// SetClipboardData 将数据写入剪切板
func SetClipboardData(data string) error {
	s, err := syscall.UTF16FromString(data)
	if err != nil {
		return fmt.Errorf("无法转换字符串: %w", err)
	}

	size := uintptr(len(s) * int(unsafe.Sizeof(s[0])))

	// 分配全局内存
	hMem, _, _ := procGlobalAlloc.Call(GMEM_MOVEABLE, size)
	if hMem == 0 {
		return fmt.Errorf("无法分配全局内存")
	}

	// 锁定内存，获取指针
	pMem, _, _ := procGlobalLock.Call(hMem)
	if pMem == 0 {
		return fmt.Errorf("无法锁定全局内存")
	}
	defer procGlobalUnlock.Call(hMem)

	// 复制数据到全局内存
	procRtlMoveMemory.Call(pMem, uintptr(unsafe.Pointer(&s[0])), size)

	// 打开剪切板
	if r, _, _ := procOpenClipboard.Call(0); r == 0 {
		return fmt.Errorf("打开剪切板失败")
	}
	defer procCloseClipboard.Call()

	// 清空剪切板
	procEmptyClipboard.Call()

	// 设置剪切板数据
	if r, _, _ := procSetClipboardData.Call(CF_UNICODETEXT, hMem); r == 0 {
		return fmt.Errorf("无法设置剪贴板数据")
	}

	return nil
}
