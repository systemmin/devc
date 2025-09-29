/**
 * @Time : 2025/3/16 10:11
 * @File : translate.go
 * @Software: dev_clip
 * @Author : Mr.Fang
 * @Description: 翻译
 */

package internal

import (
	bytes2 "bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// Translate 翻译
func Translate(api, model, content string) string {
	if len(model) != 0 {
		return TranslateOllama(api, model, content)
	}
	return TranslateYouDao(content)
}

// TranslateYouDao 翻译
func TranslateYouDao(content string) string {
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

type Input struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// TranslateOllama 翻译
func TranslateOllama(api, model, content string) string {

	// 角色设定
	sysContent := `You are a professional technical translator. 
Your task is to translate the input text between English and Simplified Chinese, depending on the input language. 
The text is mainly technical documentation, developer manuals, function/method comments, or code-related explanations. 

Translation rules:
1. Detect the input language automatically and translate it into the other language (English ↔ Simplified Chinese).
2. Keep technical terms (e.g., API, class, HTTP, JSON, SQL) in their original form unless there is a widely accepted translation.
3. Method names, function names, variable names, file paths, and code snippets must remain unchanged.
4. Use concise, formal, and professional style suitable for developer documentation.
5. Ensure the translation is accurate, consistent, and natural.
6. Do not add extra explanations, only output the translated text.`

	input := Input{
		Model: model,
		Messages: []Message{
			{
				Role:    "system",
				Content: sysContent,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: false,
	}

	client := &http.Client{}
	marshal, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/chat", api), bytes2.NewBuffer(marshal))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败 - ", err)
		return TranslateYouDao(content)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println("状态码错误 - ", resp.StatusCode)
		return TranslateYouDao(content)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	message := result["message"].(map[string]interface{})
	apiContent := message["content"].(interface{})
	return apiContent.(string)
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
