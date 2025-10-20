package main

import "testing"

func TestHello(t *testing.T) {
	assertCorrectMessage := func(t testing.TB, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}

	t.Run("in Chinese", func(t *testing.T) {
		got := Hello("Fan", "Chinese")
		want := "你好 Fan"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in English", func(t *testing.T) {
		got := Hello("Fan", "English")
		want := "Hello Fan"
		assertCorrectMessage(t, got, want)
	})

	t.Run("saying hello world when an empty string is supplied", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello World"
		assertCorrectMessage(t, got, want)
	})
}

// func TestHello(t *testing.T) {

// 	assertCorrectMessage := func(t testing.TB, got, want string) {
// 		// t.Helper 的作用是告诉测试框架这个函数是辅助函数，不要报告这个函数的行号。
// 		// 比如test不通过，不是报告t.Errorf("got %q, want %q", got, want)的行号，而是对应错误的Run里面的行号
// 		t.Helper()
// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	}

// 	t.Run("saying hello to people", func(t *testing.T) {
// 		got := Hello("Fan")
// 		want := "Hello Fan"

// 		assertCorrectMessage(t, got, want)
// 	})

// 	t.Run("saying hello world when an empty string is supplied", func(t *testing.T) {
// 		got := Hello("")
// 		want := "Hello World"

// 		assertCorrectMessage(t, got, want)
// 	})
// }

// func TestHello(t *testing.T) {

// 	t.Run("saying hello to people", func(t *testing.T) {
// 		got := Hello("Fan")
// 		want := "Hello Fan"

// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})

// 	t.Run("saying hello world when an empty string is supplied", func(t *testing.T) {
// 		got := Hello("")
// 		want := "Hello World"

// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})
// }
