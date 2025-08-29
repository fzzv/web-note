package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrNotFound = errors.New("not found")

func readConfig(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		// 包装底层错误，附加上下文
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}
	if len(b) == 0 {
		return nil, fmt.Errorf("empty config: %w", ErrNotFound)
	}
	return b, nil
}

func main() {
	_, err := readConfig("nope.json")
	if err != nil {
		fmt.Println("err:", err)
		if errors.Is(err, ErrNotFound) {
			fmt.Println("root cause: not found")
		}
	}
} 
