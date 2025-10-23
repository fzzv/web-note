package concurrency

type WebsiteChecker func(string) bool
type result struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)

	for _, url := range urls {
		// results[url] = wc(url)
		/*
			Go 中不会阻塞的操作将在称为 goroutine 的单独 进程 中运行
			要告诉 Go 开始一个新的 goroutine，
			我们把一个函数调用变成 go 声明，通过把关键字 go 放在它前面：go doSomething()。
		*/
		go func(u string) {
			/*
				当我们迭代 urls 时，不是直接写入 map，
				而是使用 send statement 将每个调用 wc 的 result 结构体发送到 resultChannel。
				这使用 <- 操作符，channel 放在左边，值放在右边
			*/
			// 发送
			resultChannel <- result{u, wc(u)}
		}(url)
	}

	for range urls {
		// 接收
		result := <-resultChannel
		results[result.string] = result.bool
	}
	return results
}
