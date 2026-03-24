package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 负责协调应用生命周期和多语言状态。
type App struct {
	counterDemo  *CounterDemo
	downloadDemo *DownloadDemo
	i18n         *Localizer
}

// NewApp 创建组合后的示例应用。
func NewApp() *App {
	i18n := NewLocalizer()

	return &App{
		counterDemo:  NewCounterDemo(),
		downloadDemo: NewDownloadDemo(i18n),
		i18n:         i18n,
	}
}

func (a *App) startup(ctx context.Context) {
	a.counterDemo.setStateFile(resolveCounterStateFile())
	a.downloadDemo.setContext(ctx)

	if err := a.counterDemo.startup(); err != nil {
		runtime.LogErrorf(ctx, "%s", a.i18n.Format("log.counter.load", err))
	}
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	selection, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         a.i18n.Text("app.close.title"),
		Message:       a.i18n.Text("app.close.message"),
		Buttons:       []string{a.i18n.Text("app.close.confirm"), a.i18n.Text("app.close.cancel")},
		DefaultButton: a.i18n.Text("app.close.cancel"),
		CancelButton:  a.i18n.Text("app.close.cancel"),
	})
	if err != nil {
		runtime.LogErrorf(ctx, "%s", a.i18n.Format("log.close.confirm", err))
		return false
	}

	return selection != a.i18n.Text("app.close.confirm")
}

func (a *App) shutdown(ctx context.Context) {
	if err := a.counterDemo.shutdown(); err != nil {
		runtime.LogErrorf(ctx, "%s", a.i18n.Format("log.counter.save", err))
	}
}

// GetLanguage 返回当前语言。
func (a *App) GetLanguage() string {
	return a.i18n.GetLanguage()
}

// SetLanguage 设置当前语言，并同步已有下载状态文案。
func (a *App) SetLanguage(language string) string {
	nextLanguage := a.i18n.SetLanguage(language)
	a.downloadDemo.SyncLanguage()
	return nextLanguage
}
