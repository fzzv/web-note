package racer

import (
	"fmt"
	"net/http"
	"time"
)

func Racer(a, b string) (winner string) {
	/*
		用time.now()来记录请求 URL 前的时间。
		然后用 http.Get 来请求 URL 的内容。
		time.Since 获取开始时间并返回一个 time.Duration 时间差
	*/
	startA := time.Now()
	http.Get(a)
	aDuration := time.Since(startA)

	startB := time.Now()
	http.Get(b)
	bDuration := time.Since(startB)

	if aDuration < bDuration {
		return a
	}
	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}

func RacerTwo(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}
	return b
}

func RacerSelect(a, b string, timeout time.Duration) (winner string, error error) {
	/*
		select construct 语句，可以帮我们轻易清晰地实现进程同步
		select 是“通道版的 switch” —— 哪个通道先准备好，select 就执行哪个分支。
		select可同时在多个 channel 上等待。
		可以使用 time.After 来防止你的系统被永久阻塞
	*/
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}

func ping(url string) chan bool {
	ch := make(chan bool)
	go func() {
		http.Get(url)
		ch <- true
	}()
	return ch
}
