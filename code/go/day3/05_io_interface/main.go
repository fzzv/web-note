package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func CopyN(r io.Reader, w io.Writer, n int64) (int64, error) {
	return io.CopyN(w, r, n)
}

func main() {
	src := strings.NewReader("Hello, IO Reader/Writer!\n")
	var dst bytes.Buffer

	if _, err := CopyN(src, &dst, 5); err != nil {
		fmt.Println("copy err:", err)
	}
	fmt.Print("dst:", dst.String())
} 
