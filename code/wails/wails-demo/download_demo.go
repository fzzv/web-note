package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const downloadProgressEventName = "demo:download:progress"

// DownloadDemo 负责文件下载进度示例。
type DownloadDemo struct {
	mu     sync.RWMutex
	ctx    context.Context
	state  DownloadState
	active bool
	i18n   *Localizer
}

func NewDownloadDemo(i18n *Localizer) *DownloadDemo {
	demo := &DownloadDemo{i18n: i18n}
	demo.state = demo.newIdleState()

	return demo
}

func (d *DownloadDemo) setContext(ctx context.Context) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.ctx = ctx
}

func (d *DownloadDemo) GetState() DownloadState {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.state
}

func (d *DownloadDemo) Reset() (DownloadState, error) {
	d.mu.Lock()
	if d.active {
		state := d.state
		d.mu.Unlock()
		return state, errors.New(d.i18n.Text("download.reset.busy"))
	}

	state := d.newIdleState()
	ctx := d.ctx
	d.state = state
	d.mu.Unlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, downloadProgressEventName, state)
	}

	return state, nil
}

func (d *DownloadDemo) StartDownload(rawURL string) (DownloadState, error) {
	trimmedURL := strings.TrimSpace(rawURL)
	if trimmedURL == "" {
		return d.GetState(), errors.New(d.i18n.Text("download.url.empty"))
	}

	targetURL, err := neturl.ParseRequestURI(trimmedURL)
	if err != nil {
		return d.GetState(), fmt.Errorf(d.i18n.Text("download.url.invalid"), err)
	}

	fileName := fileNameFromURL(targetURL)
	destinationDir := resolveDownloadDirectory()
	if err := os.MkdirAll(destinationDir, 0o755); err != nil {
		return d.GetState(), fmt.Errorf(d.i18n.Text("download.dir.create"), err)
	}

	destination, err := nextAvailableFilePath(destinationDir, fileName)
	if err != nil {
		return d.GetState(), fmt.Errorf(d.i18n.Text("download.path.prepare"), err)
	}

	state := DownloadState{
		Status:      "starting",
		Message:     d.i18n.Text("download.starting"),
		URL:         trimmedURL,
		FileName:    filepath.Base(destination),
		Destination: destination,
		Progress:    0,
	}

	d.mu.Lock()
	if d.active {
		current := d.state
		d.mu.Unlock()
		return current, errors.New(d.i18n.Text("download.busy"))
	}

	ctx := d.ctx
	d.state = state
	d.active = true
	d.mu.Unlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, downloadProgressEventName, state)
	}

	go d.runDownload(trimmedURL, destination)

	return state, nil
}

func (d *DownloadDemo) runDownload(rawURL, destination string) {
	fileName := filepath.Base(destination)

	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		d.finishWithError(rawURL, destination, fileName, 0, -1, fmt.Errorf(d.i18n.Text("download.request.build"), err))
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		d.finishWithError(rawURL, destination, fileName, 0, -1, fmt.Errorf(d.i18n.Text("download.request.send"), err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		d.finishWithError(
			rawURL,
			destination,
			fileName,
			0,
			resp.ContentLength,
			fmt.Errorf(d.i18n.Text("download.response.status"), resp.Status),
		)
		return
	}

	file, err := os.Create(destination)
	if err != nil {
		d.finishWithError(
			rawURL,
			destination,
			fileName,
			0,
			resp.ContentLength,
			fmt.Errorf(d.i18n.Text("download.file.create"), err),
		)
		return
	}
	defer file.Close()

	totalBytes := resp.ContentLength
	d.publish(DownloadState{
		Status:      "downloading",
		Message:     d.i18n.Text("download.running"),
		URL:         rawURL,
		FileName:    fileName,
		Destination: destination,
		TotalBytes:  totalBytes,
		Progress:    0,
	}, true)

	buffer := make([]byte, 32*1024)
	var downloadedBytes int64
	lastEmitAt := time.Now()

	for {
		n, readErr := resp.Body.Read(buffer)
		if n > 0 {
			if _, err := file.Write(buffer[:n]); err != nil {
				d.finishWithError(
					rawURL,
					destination,
					fileName,
					downloadedBytes,
					totalBytes,
					fmt.Errorf(d.i18n.Text("download.file.write"), err),
				)
				return
			}

			downloadedBytes += int64(n)
			if shouldEmitProgress(totalBytes, downloadedBytes, lastEmitAt) {
				d.publish(DownloadState{
					Status:          "downloading",
					Message:         d.progressMessage(downloadedBytes, totalBytes),
					URL:             rawURL,
					FileName:        fileName,
					Destination:     destination,
					DownloadedBytes: downloadedBytes,
					TotalBytes:      totalBytes,
					Progress:        calculateProgress(downloadedBytes, totalBytes),
				}, true)
				lastEmitAt = time.Now()
			}
		}

		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			d.finishWithError(
				rawURL,
				destination,
				fileName,
				downloadedBytes,
				totalBytes,
				fmt.Errorf(d.i18n.Text("download.response.read"), readErr),
			)
			return
		}
	}

	d.publish(DownloadState{
		Status:          "completed",
		Message:         d.i18n.Format("download.complete", formatBytes(downloadedBytes), destination),
		URL:             rawURL,
		FileName:        fileName,
		Destination:     destination,
		DownloadedBytes: downloadedBytes,
		TotalBytes:      downloadedBytes,
		Progress:        100,
	}, false)
}

// SyncLanguage 用当前语言刷新已有状态文案。
func (d *DownloadDemo) SyncLanguage() {
	d.mu.RLock()
	state := d.state
	ctx := d.ctx
	d.mu.RUnlock()

	nextState := d.localizeState(state)
	if nextState.Message == state.Message {
		return
	}

	d.mu.Lock()
	d.state = nextState
	d.mu.Unlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, downloadProgressEventName, nextState)
	}
}

func (d *DownloadDemo) finishWithError(rawURL, destination, fileName string, downloadedBytes, totalBytes int64, err error) {
	_ = os.Remove(destination)

	d.publish(DownloadState{
		Status:          "error",
		Message:         err.Error(),
		URL:             rawURL,
		FileName:        fileName,
		Destination:     destination,
		DownloadedBytes: downloadedBytes,
		TotalBytes:      totalBytes,
		Progress:        calculateProgress(downloadedBytes, totalBytes),
	}, false)
}

func (d *DownloadDemo) publish(state DownloadState, active bool) {
	d.mu.Lock()
	ctx := d.ctx
	d.state = state
	d.active = active
	d.mu.Unlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, downloadProgressEventName, state)
	}
}

func (d *DownloadDemo) newIdleState() DownloadState {
	return DownloadState{
		Status:  "idle",
		Message: d.i18n.Text("download.idle"),
	}
}

func (d *DownloadDemo) localizeState(state DownloadState) DownloadState {
	switch state.Status {
	case "idle":
		state.Message = d.i18n.Text("download.idle")
	case "starting":
		state.Message = d.i18n.Text("download.starting")
	case "downloading":
		if state.DownloadedBytes > 0 {
			state.Message = d.progressMessage(state.DownloadedBytes, state.TotalBytes)
		} else {
			state.Message = d.i18n.Text("download.running")
		}
	case "completed":
		state.Message = d.i18n.Format("download.complete", formatBytes(state.DownloadedBytes), state.Destination)
	}

	return state
}

func fileNameFromURL(targetURL *neturl.URL) string {
	fileName := path.Base(targetURL.Path)
	if fileName == "." || fileName == "/" || fileName == "" {
		fileName = fmt.Sprintf("download-%s.bin", time.Now().Format("20060102-150405"))
	}

	return fileName
}

func resolveDownloadDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "wails-demo-downloads")
	}

	return filepath.Join(homeDir, "Downloads", "wails-demo")
}

func nextAvailableFilePath(dir, fileName string) (string, error) {
	ext := filepath.Ext(fileName)
	base := strings.TrimSuffix(fileName, ext)
	if base == "" {
		base = "download"
	}

	candidate := filepath.Join(dir, fileName)
	index := 1

	for {
		_, err := os.Stat(candidate)
		if errors.Is(err, os.ErrNotExist) {
			return candidate, nil
		}
		if err != nil {
			return "", err
		}

		candidate = filepath.Join(dir, fmt.Sprintf("%s-%d%s", base, index, ext))
		index++
	}
}

func shouldEmitProgress(totalBytes, downloadedBytes int64, lastEmitAt time.Time) bool {
	if time.Since(lastEmitAt) >= 150*time.Millisecond {
		return true
	}

	return totalBytes > 0 && downloadedBytes >= totalBytes
}

func calculateProgress(downloadedBytes, totalBytes int64) float64 {
	if totalBytes <= 0 {
		return 0
	}

	return float64(downloadedBytes) / float64(totalBytes) * 100
}

func (d *DownloadDemo) progressMessage(downloadedBytes, totalBytes int64) string {
	if totalBytes <= 0 {
		return d.i18n.Format("download.progress.partial", formatBytes(downloadedBytes))
	}

	return d.i18n.Format("download.progress.total", formatBytes(downloadedBytes), formatBytes(totalBytes))
}

func formatBytes(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	units := []string{"KB", "MB", "GB", "TB"}
	value := float64(size) / 1024
	unit := 0

	for value >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}

	return fmt.Sprintf("%.1f %s", value, units[unit])
}
