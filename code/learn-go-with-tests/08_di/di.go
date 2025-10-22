package main

import (
	"fmt"
	"io"
	"net/http"
)

/*
文件是 Writer；
网络连接（socket）是 Writer；
内存缓冲区（如 bytes.Buffer）也是 Writer；
甚至你自己定义的结构体也能成为 Writer。
*/

func Greet(w io.Writer, name string) {
	// fmt.Fprintf 和 fmt.Printf 一样
	// 只不过 fmt.Fprintf 会接收一个 Writer 参数，用于把字符串传递过去
	// 而 fmt.Printf 默认是标准输出
	fmt.Fprintf(w, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	Greet(w, "Fan")
}

func main() {
	http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler))
}
