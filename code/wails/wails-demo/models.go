package main

// CounterState 表示计数器持久化到本地的内容。
type CounterState struct {
	Count int `json:"count"`
}

// DownloadState 表示下载示例推送给前端的状态。
type DownloadState struct {
	Status          string  `json:"status"`
	Message         string  `json:"message"`
	URL             string  `json:"url"`
	FileName        string  `json:"fileName"`
	Destination     string  `json:"destination"`
	DownloadedBytes int64   `json:"downloadedBytes"`
	TotalBytes      int64   `json:"totalBytes"`
	Progress        float64 `json:"progress"`
}
