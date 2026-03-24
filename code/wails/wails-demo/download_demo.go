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

const downloadIdleMessage = "Enter a URL to stream a file and watch backend progress events in real time."

// DownloadDemo demonstrates progress notifications from Go to the frontend.
type DownloadDemo struct {
	mu     sync.RWMutex
	ctx    context.Context
	state  DownloadState
	active bool
}

func NewDownloadDemo() *DownloadDemo {
	return &DownloadDemo{
		state: DownloadState{
			Status:  "idle",
			Message: downloadIdleMessage,
		},
	}
}

func (d *DownloadDemo) SetContext(ctx context.Context) {
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
		return state, errors.New("a download is still in progress")
	}

	state := DownloadState{
		Status:  "idle",
		Message: downloadIdleMessage,
	}
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
		return d.GetState(), errors.New("download url cannot be empty")
	}

	targetURL, err := neturl.ParseRequestURI(trimmedURL)
	if err != nil {
		return d.GetState(), fmt.Errorf("invalid download url: %w", err)
	}

	fileName := fileNameFromURL(targetURL)
	destinationDir := resolveDownloadDirectory()
	if err := os.MkdirAll(destinationDir, 0o755); err != nil {
		return d.GetState(), fmt.Errorf("create download directory: %w", err)
	}

	destination, err := nextAvailableFilePath(destinationDir, fileName)
	if err != nil {
		return d.GetState(), fmt.Errorf("prepare download path: %w", err)
	}

	state := DownloadState{
		Status:      "starting",
		Message:     "Preparing request and waiting for response headers...",
		URL:         trimmedURL,
		FileName:    filepath.Base(destination),
		Destination: destination,
		Progress:    0,
	}

	d.mu.Lock()
	if d.active {
		current := d.state
		d.mu.Unlock()
		return current, errors.New("a download is already running")
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
		d.finishWithError(rawURL, destination, fileName, 0, -1, fmt.Errorf("build request: %w", err))
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		d.finishWithError(rawURL, destination, fileName, 0, -1, fmt.Errorf("send request: %w", err))
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
			fmt.Errorf("unexpected response status: %s", resp.Status),
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
			fmt.Errorf("create local file: %w", err),
		)
		return
	}
	defer file.Close()

	totalBytes := resp.ContentLength
	d.publish(DownloadState{
		Status:      "downloading",
		Message:     "Response received. Streaming file contents from Go to disk...",
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
					fmt.Errorf("write local file: %w", err),
				)
				return
			}

			downloadedBytes += int64(n)
			if shouldEmitProgress(totalBytes, downloadedBytes, lastEmitAt) {
				d.publish(DownloadState{
					Status:          "downloading",
					Message:         progressMessage(downloadedBytes, totalBytes),
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
				fmt.Errorf("read response body: %w", readErr),
			)
			return
		}
	}

	d.publish(DownloadState{
		Status:          "completed",
		Message:         fmt.Sprintf("Download complete. Saved %s to %s", formatBytes(downloadedBytes), destination),
		URL:             rawURL,
		FileName:        fileName,
		Destination:     destination,
		DownloadedBytes: downloadedBytes,
		TotalBytes:      downloadedBytes,
		Progress:        100,
	}, false)
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

func progressMessage(downloadedBytes, totalBytes int64) string {
	if totalBytes <= 0 {
		return fmt.Sprintf("Downloaded %s", formatBytes(downloadedBytes))
	}

	return fmt.Sprintf("Downloaded %s of %s", formatBytes(downloadedBytes), formatBytes(totalBytes))
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
