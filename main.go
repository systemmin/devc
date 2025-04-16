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
var top bool
var han int
var hump int
var tran bool
var w bool

var dicts []string

var devPath = "./"

func init() {
	flag.BoolVar(&top, "top", false, "[窗口] 置顶当前 CMD 窗口")
	flag.IntVar(&han, "han", 4, "[关键字翻译] 汉语关键字长度，默认支持4个汉字，最少2个汉字。例如：订单编号、订单状态、创建时间")
	flag.IntVar(&hump, "hump", 1, "[关键字翻译] 将翻译内容写入剪切板，命名规则：1.大驼峰 2.小驼峰 3.下划线")
	flag.BoolVar(&tran, "tran", false, "[翻译工具] 开启自动检测翻译，关键字翻译将关闭")
	flag.BoolVar(&w, "w", false, "[翻译工具|关键字翻译] 将翻译内容写入剪切板")
	// 解析命令行参数
	flag.Parse()

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
	if top {
		internal.TopUp()
	}
	length := han
	if length < 2 {
		length = 2
	}

	// 临时写入
	tempWrite := ""
	tempData := ""
	// 开始监听
	watch := internal.Watch()
	for data := range watch {
		// 跳过空内容、跳过写入剪切板内容
		if len(strings.TrimSpace(data)) == 0 || strings.Contains(tempWrite, data) || tempWrite == data || strings.Contains(tempData, data) || tempData == data {
			continue
		}
		// 历史记录
		go writeHistory(data)
		tempData = data

		if tran {
			translate := internal.Translate(data)
			fmt.Println("原文->", data)
			fmt.Println("翻译->", translate)
			if w && len(translate) > 0 {
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
					if w {
						tempWrite = internal.SetData(line, hump)
					}
				} else if len(data) <= length*3 {
					translate := internal.Translate(data)
					if len(translate) > 0 {
						line = fmt.Sprintf("%s|%s\n", data, internal.ToCamelCase(translate))
						fmt.Print("翻译->" + line)
						go internal.SaveDicts(devPath, line)
						dicts = append(dicts, line)
						if w {
							tempWrite = internal.SetData(line, hump)
						}
					}
				}
			} else if textType == internal.ENGLISH { // 英文
				line := internal.FindDict(data, dicts)
				if len(line) > 0 {
					fmt.Println("缓存->" + line)
					if w {
						tempWrite = internal.SetData(line, hump)
					}
				} else {
					translate := internal.Translate(data)
					fmt.Println("英文->", data)
					fmt.Println("中文->", translate)
				}
			} else if textType == internal.CAMEL_CASE { // 驼峰
				word := internal.CamelToWord(data)
				line := internal.FindDict(data, dicts)
				fmt.Println("格式转换->" + word)
				if len(line) > 0 {
					fmt.Println("字典缓存->" + strings.Split(line, "|")[0])
				} else {
					translate := internal.Translate(word)
					fmt.Println("转换翻译->", translate)
				}
			} else {
				fmt.Println("排除内容->" + data)
			}
		}
	}
}
