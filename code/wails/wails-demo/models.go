package main

// CounterState is the persisted counter payload.
type CounterState struct {
	Count int `json:"count"`
}

// DownloadState is emitted to the frontend during the download demo.
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
