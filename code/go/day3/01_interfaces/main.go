package main

import "fmt"

type Notifier interface {
	Notify(msg string) error
}

type Email struct{ Addr string }
func (e Email) Notify(msg string) error {
	fmt.Printf("send email to %s: %s\n", e.Addr, msg)
	return nil
}

type SMS struct{ Number string }
func (s SMS) Notify(msg string) error {
	fmt.Printf("send sms to %s: %s\n", s.Number, msg)
	return nil
}

// 依赖倒置：依赖接口而非实现
func Alert(n Notifier, content string) error { return n.Notify(content) }

func main() {
	_ = Alert(Email{Addr: "a@example.com"}, "hello via email")
	_ = Alert(SMS{Number: "+86-123456"}, "hello via sms")
} 
