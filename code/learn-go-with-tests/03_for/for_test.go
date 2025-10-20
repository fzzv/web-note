package iteration

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {
	repeated := Repeat("a", 1)
	expected := "a"

	if repeated != expected {
		t.Errorf("expected '%q' but got '%q'", expected, repeated)
	}
}

// 使用 go test -bench="." 来运行基准测试
func BenchmarkRepeat(b *testing.B) {
	// https://go.dev/blog/testing-b-loop
	// 使用 b.Loop() 代替 range b.N
	for b.Loop() {
		Repeat("a", 1)
	}
}

func ExampleRepeat() {
	repeated := Repeat("a", 1)
	fmt.Println(repeated)
	// Output: a
}
