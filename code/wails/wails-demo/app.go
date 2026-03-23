package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	mu        sync.RWMutex
	count     int
	stateFile string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.stateFile = resolveStateFile()

	if err := a.loadCount(); err != nil {
		runtime.LogErrorf(ctx, "load counter state failed: %v", err)
	}
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	selection, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Close counter?",
		Message:       "Do you want to close the app? The current count will be saved on shutdown.",
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "No",
		CancelButton:  "No",
	})
	if err != nil {
		runtime.LogErrorf(ctx, "show close confirmation failed: %v", err)
		return false
	}

	return selection != "Yes"
}

func (a *App) shutdown(ctx context.Context) {
	if err := a.saveCount(); err != nil {
		runtime.LogErrorf(ctx, "save counter state failed: %v", err)
	}
}

func (a *App) Increment() int {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.count++
	return a.count
}

func (a *App) Decrement() int {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.count--
	return a.count
}

func (a *App) GetCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.count
}

func (a *App) loadCount() error {
	data, err := os.ReadFile(a.stateFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	var state CounterState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	a.mu.Lock()
	a.count = state.Count
	a.mu.Unlock()

	return nil
}

func (a *App) saveCount() error {
	a.mu.RLock()
	state := CounterState{Count: a.count}
	a.mu.RUnlock()

	if err := os.MkdirAll(filepath.Dir(a.stateFile), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(a.stateFile, data, 0o644)
}

func resolveStateFile() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "counter-state.json"
	}

	return filepath.Join(configDir, "wails-demo", "counter-state.json")
}
