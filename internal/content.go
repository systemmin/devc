/**
 * @Time : 2025/3/16 10:26
 * @File : content.go
 * @Software: dev_clip
 * @Author : Mr.Fang
 * @Description: 剪切板内容处理
 */

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ToCamelCase 驼峰命名，返回内容 Chinese character|ChineseCharacter|chineseCharacter|chinese_character
func ToCamelCase(content string) string {
	fields := strings.Split(content, " ")
	var greatHump []string
	var smallHump []string
	var underLine []string
	var result []string
	for i, field := range fields {
		title := strings.Title(field)
		lower := strings.ToLower(field)
		greatHump = append(greatHump, title)
		if i == 0 {
			// 第一个单词首字母小写
			smallHump = append(smallHump, lower)
		} else {
			smallHump = append(smallHump, title)
		}
		underLine = append(underLine, lower)
	}
	result = append(result, content, strings.Join(greatHump, ""), strings.Join(smallHump, ""), strings.Join(underLine, "_"))
	return strings.Join(result, "|")
}

// AppendFile 追加写入
func AppendFile(path, content string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("打开文件失败: %v\n", err)
		return
	}
	defer file.Close()
	// 写入内容
	if _, err := file.WriteString(content); err != nil {
		fmt.Printf("写入文件失败: %v\n", err)
		return
	}
}

// FindDict 查找已缓存内容
func FindDict(prefix string, dicts []string) string {
	for _, line := range dicts {
		split := strings.Split(line, "|")
		for _, block := range split {
			if block == prefix {
				return line
			}
		}
	}
	return ""
}

// GetDicts 获取缓存字典
func GetDicts(devPath string) []string {
	file, err := os.ReadFile(filepath.Join(devPath, "dict.txt"))
	if err != nil {
		fmt.Println("暂无字典数据")
		return nil
	}
	return strings.Split(string(file), "\n")
}

// SaveDicts 获取缓存字典
func SaveDicts(devPath, content string) {
	AppendFile(filepath.Join(devPath, "dict.txt"), content)
}

// SetData 写入剪切板
// 0 1 2 3 4
// 0 0 1 2 3
func SetData(line string, index int) string {
	if index == 0 {
		return ""
	}
	content := ""
	split := strings.Split(line, "|")
	content = split[index+1]
	err := SetClipboardData(content)
	if err != nil {
		fmt.Println("写入失败", err)
	}
	return content
}
