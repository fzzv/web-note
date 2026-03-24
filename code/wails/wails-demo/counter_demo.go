package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// CounterDemo 负责持久化计数器示例。
type CounterDemo struct {
	mu        sync.RWMutex
	count     int
	stateFile string
}

func NewCounterDemo() *CounterDemo {
	return &CounterDemo{}
}

func (d *CounterDemo) setStateFile(path string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.stateFile = path
}

func (d *CounterDemo) startup() error {
	return d.load()
}

func (d *CounterDemo) shutdown() error {
	return d.save()
}

func (d *CounterDemo) Increment() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.count++
	return d.count
}

func (d *CounterDemo) Decrement() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.count--
	return d.count
}

func (d *CounterDemo) GetCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.count
}

func (d *CounterDemo) load() error {
	d.mu.RLock()
	stateFile := d.stateFile
	d.mu.RUnlock()

	data, err := os.ReadFile(stateFile)
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

	d.mu.Lock()
	d.count = state.Count
	d.mu.Unlock()

	return nil
}

func (d *CounterDemo) save() error {
	d.mu.RLock()
	state := CounterState{Count: d.count}
	stateFile := d.stateFile
	d.mu.RUnlock()

	if err := os.MkdirAll(filepath.Dir(stateFile), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(stateFile, data, 0o644)
}

func resolveCounterStateFile() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "counter-state.json"
	}

	return filepath.Join(configDir, "wails-demo", "counter-state.json")
}
