/**
 * @Time : 2025/3/18 10:57
 * @File : classify_text.go
 * @Software: dev_clip
 * @Author : Mr.Fang
 * @Description: 文本内容分类
 */

package internal

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	CHINESE = iota // 中文
	CAMEL_CASE
	SNAKE_CASE
	CODE_SNIPPET
	ENGLISH
	OTHER
)

// 纯中文检测
func isChinese(text string) bool {
	re := regexp.MustCompile(`^[\p{Han}]+$`)
	return re.MatchString(text)
}

// IsCamelCase 驼峰命名检测
func IsCamelCase(text string) bool {
	// 小驼峰
	re := regexp.MustCompile(`^[a-z]+(?:[A-Z][a-z]*)+$`)
	match := re.MatchString(text)
	if !match {
		// 大驼峰
		re = regexp.MustCompile(`^[A-Z][a-z]+(?:[A-Z][a-z]*)+$`)
		match = re.MatchString(text)
	}
	if !match {
		// 下划线
		match = isSnakeCase(text)
	}
	return match
}

// 下划线命名检测
func isSnakeCase(text string) bool {
	re := regexp.MustCompile(`^[a-z]+(?:_[a-z]+)+$`)
	return re.MatchString(text)
}

// IsCodeSnippet 代码片段检测（通过特殊符号判断）
func IsCodeSnippet(text string) bool {
	codeChars := []string{"{", "}", ";", "(", ")", "[", "]", "<", ">", "=", "\""}
	for _, char := range codeChars {
		if strings.Contains(text, char) {
			// [A-Za-z]+\(
			if char == "(" {
				re := regexp.MustCompile(`[A-Za-z]+\(+`)
				return re.MatchString(text)
			}
			return true
		}
	}
	// 第二种情况；空格拆分只有一个长度，并且中间有.连接
	if len(strings.Split(text, " ")) == 1 && strings.Contains(text, ".") {
		_, b := ExtractURLs(text)
		if b {
			return false
		}
		return true
	}
	return false
}

// IsEnglish 纯英文检测（排除驼峰命名、下划线命名、代码片段）
func IsEnglish(text string) bool {
	if IsCamelCase(text) || isSnakeCase(text) || IsCodeSnippet(text) {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z\s.,!?'’]+$`)
	if re.MatchString(text) {
		return true
	}
	return IsEnglishContent(text)
}

// IsEnglishContent 内容检测，只要英文占比超过 60%
func IsEnglishContent(text string) bool {
	re := regexp.MustCompile(`^[a-zA-Z\s.,!?'’]+$`)
	split := strings.Split(text, " ")
	l := len(split)
	count := 0
	for _, s := range split {
		if re.MatchString(s) {
			count++
		}
	}
	percentage, _ := CalcPercentage(float64(count), float64(l))
	return percentage > 60
}

// CalcPercentage 计算百分比 numerator 分母 denominator 分子
func CalcPercentage(numerator, denominator float64) (int, error) {
	if denominator == 0 {
		return 0, fmt.Errorf("除数为 0 ")
	}
	percentage := (numerator / denominator) * 100
	return int(percentage), nil
}

// ExtractURLs 提取 URL，并检测是否全是 URL
func ExtractURLs(text string) ([]string, bool) {
	re := regexp.MustCompile(`https?://[^\s]+`)
	urls := re.FindAllString(text, -1)
	return urls, len(urls) > 0 && len(urls[0]) == len(text)
}

// ClassifyText 分类文本
func ClassifyText(text string) int {
	if isChinese(text) {
		fmt.Println("纯中文内容")
		return CHINESE
	} else if IsCamelCase(text) {
		fmt.Println("驼峰命名")
		return CAMEL_CASE
	} else if IsCodeSnippet(text) {
		fmt.Println("代码片段")
		return CODE_SNIPPET
	} else if IsEnglish(text) {
		fmt.Println("纯英文内容")
		return ENGLISH
	} else {
		fmt.Println("无法分类的内容")
		return OTHER
	}
}
