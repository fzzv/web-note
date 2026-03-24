package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App wires lifecycle hooks to individual demos.
type App struct {
	counterDemo  *CounterDemo
	downloadDemo *DownloadDemo
}

// NewApp creates a composed demo application.
func NewApp() *App {
	return &App{
		counterDemo:  NewCounterDemo(),
		downloadDemo: NewDownloadDemo(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.counterDemo.SetStateFile(resolveCounterStateFile())
	a.downloadDemo.SetContext(ctx)

	if err := a.counterDemo.Startup(); err != nil {
		runtime.LogErrorf(ctx, "load counter state failed: %v", err)
	}
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	selection, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Close demos?",
		Message:       "Do you want to close the demo workspace? Counter state will be saved on shutdown.",
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
	if err := a.counterDemo.Shutdown(); err != nil {
		runtime.LogErrorf(ctx, "save counter state failed: %v", err)
	}
}
