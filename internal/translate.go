/**
 * @Time : 2025/3/16 10:11
 * @File : translate.go
 * @Software: dev_clip
 * @Author : Mr.Fang
 * @Description: 翻译
 */

package internal

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Translate 翻译
func Translate(content string) string {
	client := &http.Client{}
	var data = strings.NewReader("inputtext=" + content + "&type=AUTO")
	req, err := http.NewRequest("POST", "https://m.youdao.com/translate", data)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://m.youdao.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://m.youdao.com/translate")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	result := string(bytes)
	return matchContent(result)
}

// 匹配内容
func matchContent(t string) string {
	compile := regexp.MustCompile(`<ul id="translateResult">([\s\S]*?)<\/ul>`)
	matches := compile.FindStringSubmatch(t)
	if len(matches) > 1 {
		content := matches[1]
		content = strings.TrimSpace(content)
		content = strings.ReplaceAll(content, "<li>", "")
		content = strings.ReplaceAll(content, "</li>", "")
		return content
	}
	return ""
}
