package integers

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4
	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum)
	}
}

// Example方法可以在 godoc 的example中显示
func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// 下面这个注释删除了，函数会被编译，但不会被执行
	// Output: 6
}
