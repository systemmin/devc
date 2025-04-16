/**
 * @Time : 2025/4/15 17:16
 * @File : main.go
 * @Software: devc
 * @Author : Mr.Fang
 * @Description:
 */

package main

import (
	"devc/internal"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 中文字符长度，需要乘 3，中文占 3 个字节
// ==================== 命令行参数
var top = flag.Bool("top", false, "[窗口] 置顶当前 CMD 窗口")
var han = flag.Int("han", 4, "[关键字翻译] 汉语关键字长度，默认支持4个汉字，最少2个汉字。例如：订单编号、订单状态、创建时间")
var hump = flag.Int("hump", 1, "[关键字翻译] 将翻译内容写入剪切板，命名规则：1.大驼峰 2.小驼峰 3.下划线")
var tran = flag.Bool("tran", false, "[翻译工具] 开启自动检测翻译，关键字翻译将关闭")
var w = flag.Bool("w", false, "[翻译工具] 将翻译内容写入剪切板")

var dicts []string

var devPath = "./"

func init() {
	// 解析命令行参数
	flag.Parse()

	// 打印参数值
	//fmt.Println("top:", *top)
	//fmt.Println("han:", *han)
	//fmt.Println("hump:", *hump)
	//fmt.Println("tran:", *tran)

	// 应用程序数据的存储目录
	devPath = os.Getenv("DEVC")
	if len(devPath) == 0 {
		devPath = filepath.Join(os.Getenv("APPDATA"), "devc") // 获取应用数据目录 %APPDATA%\devc
		if _, err := os.Stat(devPath); err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(devPath, os.ModePerm)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	// ==================== 缓存数据
	dicts = internal.GetDicts(devPath)
}

func writeHistory(content string) {
	format := time.Now().Format(time.DateOnly)
	internal.AppendFile(filepath.Join(devPath, fmt.Sprintf("./%s_history.txt", format)), content+"\n")
}

func main() {
	// 置顶窗口
	if *top {
		internal.TopUp()
	}
	length := *han
	if length < 2 {
		length = 2
	}
	// 临时写入
	tempWrite := ""
	// 开始监听
	watch := internal.Watch()
	for data := range watch {
		// 跳过空内容、跳过写入剪切板内容
		if len(strings.TrimSpace(data)) == 0 || strings.Contains(tempWrite, data) || tempWrite == data {
			continue
		}

		// 历史记录
		go writeHistory(data)

		if *tran {
			translate := internal.Translate(data)
			fmt.Println("原文->", data)
			fmt.Println("翻译->", translate)
			if *w && len(translate) > 0 {
				err := internal.SetClipboardData(translate)
				if err != nil {
					fmt.Println("写入失败")
				}
				tempWrite = translate
			}
		} else {
			// 内容分类
			textType := internal.ClassifyText(data)
			if textType == internal.CHINESE { // 中文
				line := internal.FindDict(data, dicts)
				if len(line) > 0 {
					fmt.Println("缓存->" + line)
					tempWrite = internal.SetData(line, *hump)
				} else if len(data) <= length*3 {
					translate := internal.Translate(data)
					if len(translate) > 0 {
						line = fmt.Sprintf("%s|%s\n", data, internal.ToCamelCase(translate))
						fmt.Print("翻译->" + line)
						go internal.SaveDicts(devPath, line)
						dicts = append(dicts, line)
						tempWrite = internal.SetData(line, *hump)
					}
				}
			} else if textType == internal.ENGLISH { // 英文
				line := internal.FindDict(data, dicts)
				if len(line) > 0 {
					fmt.Println("缓存->" + line)
					tempWrite = internal.SetData(line, *hump)
				} else {
					translate := internal.Translate(data)
					fmt.Println("中文->", data)
					fmt.Println("英文->", translate)
				}
			} else {
				fmt.Println("排除内容->" + data)
			}
		}
	}
}
