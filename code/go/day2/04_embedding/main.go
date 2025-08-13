package main

import "fmt"

type Logger struct{}
func (Logger) Log(msg string) { fmt.Println("[LOG]", msg) }

type Service struct {
	Logger // 组合：直接提升方法
	Name   string
}

func (s Service) Work() {
	s.Log("service:" + s.Name)
}

func main() {
	s := Service{Name: "Note"}
	s.Log("service Log:" + s.Name)
	s.Work() // 直接调用提升的方法
} 
