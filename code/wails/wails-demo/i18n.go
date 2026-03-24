package main

import (
	"fmt"
	"strings"
	"sync"
)

// Language 表示界面语言代码。
type Language string

const (
	LanguageChinese Language = "zh-CN"
	LanguageEnglish Language = "en-US"
)

var i18nMessages = map[Language]map[string]string{
	LanguageChinese: {
		"app.close.title":           "关闭示例集？",
		"app.close.message":         "确认关闭示例工作台吗？当前计数会在退出时自动保存。",
		"app.close.confirm":         "是",
		"app.close.cancel":          "否",
		"log.counter.load":          "加载计数器状态失败: %v",
		"log.counter.save":          "保存计数器状态失败: %v",
		"log.close.confirm":         "显示关闭确认框失败: %v",
		"download.idle":             "输入一个文件地址，观察后端实时推送下载进度事件。",
		"download.reset.busy":       "当前仍有下载任务在进行中",
		"download.url.empty":        "下载地址不能为空",
		"download.url.invalid":      "下载地址格式不正确: %v",
		"download.dir.create":       "创建下载目录失败: %v",
		"download.path.prepare":     "准备下载文件路径失败: %v",
		"download.starting":         "正在准备请求并等待响应头...",
		"download.running":          "已收到响应，Go 正在把文件内容写入本地磁盘...",
		"download.busy":             "当前已有下载任务正在执行",
		"download.request.build":    "构建下载请求失败: %v",
		"download.request.send":     "发送下载请求失败: %v",
		"download.response.status":  "响应状态异常: %s",
		"download.file.create":      "创建本地文件失败: %v",
		"download.file.write":       "写入本地文件失败: %v",
		"download.response.read":    "读取响应内容失败: %v",
		"download.complete":         "下载完成，已将 %s 保存到 %s",
		"download.progress.partial": "已下载 %s",
		"download.progress.total":   "已下载 %s / %s",
	},
	LanguageEnglish: {
		"app.close.title":           "Close demos?",
		"app.close.message":         "Do you want to close the demo workspace? Counter state will be saved on shutdown.",
		"app.close.confirm":         "Yes",
		"app.close.cancel":          "No",
		"log.counter.load":          "load counter state failed: %v",
		"log.counter.save":          "save counter state failed: %v",
		"log.close.confirm":         "show close confirmation failed: %v",
		"download.idle":             "Enter a URL to stream a file and watch backend progress events in real time.",
		"download.reset.busy":       "a download is still in progress",
		"download.url.empty":        "download URL cannot be empty",
		"download.url.invalid":      "invalid download URL: %v",
		"download.dir.create":       "create download directory: %v",
		"download.path.prepare":     "prepare download file path: %v",
		"download.starting":         "Preparing request and waiting for response headers...",
		"download.running":          "Response received. Streaming file contents from Go to disk...",
		"download.busy":             "a download is already running",
		"download.request.build":    "build request: %v",
		"download.request.send":     "send request: %v",
		"download.response.status":  "unexpected response status: %s",
		"download.file.create":      "create local file: %v",
		"download.file.write":       "write local file: %v",
		"download.response.read":    "read response body: %v",
		"download.complete":         "Download complete. Saved %s to %s",
		"download.progress.partial": "Downloaded %s",
		"download.progress.total":   "Downloaded %s / %s",
	},
}

// Localizer 保存当前语言，并负责返回对应文案。
type Localizer struct {
	mu       sync.RWMutex
	language Language
}

// NewLocalizer 创建多语言管理器。
func NewLocalizer() *Localizer {
	return &Localizer{
		language: LanguageChinese,
	}
}

// GetLanguage 返回当前语言代码。
func (l *Localizer) GetLanguage() string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return string(l.language)
}

// SetLanguage 更新当前语言代码，并返回规范化后的结果。
func (l *Localizer) SetLanguage(language string) string {
	nextLanguage := normalizeLanguage(language)

	l.mu.Lock()
	l.language = nextLanguage
	l.mu.Unlock()

	return string(nextLanguage)
}

// Text 返回指定键对应的文案。
func (l *Localizer) Text(key string) string {
	l.mu.RLock()
	language := l.language
	l.mu.RUnlock()

	if text, ok := i18nMessages[language][key]; ok {
		return text
	}

	return i18nMessages[LanguageChinese][key]
}

// Format 返回格式化后的多语言文案。
func (l *Localizer) Format(key string, args ...any) string {
	return fmt.Sprintf(l.Text(key), args...)
}

func normalizeLanguage(language string) Language {
	normalized := strings.ToLower(strings.TrimSpace(language))
	switch {
	case strings.HasPrefix(normalized, "en"):
		return LanguageEnglish
	case strings.HasPrefix(normalized, "zh"):
		return LanguageChinese
	default:
		return LanguageChinese
	}
}
