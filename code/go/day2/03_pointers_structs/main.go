package main

import "fmt"

type User struct {
	ID   int
	Name string
}

func NewUser(id int, name string) *User {
	return &User{ID: id, Name: name}
}

func (u *User) Rename(name string) { u.Name = name }

func main() {
	u := NewUser(1, "Alice")
	fmt.Println("before:", *u)
	u.Rename("Alicia")
	fmt.Println("after:", *u)

	// 指针与值拷贝
	copy := *u
	copy.Rename("Copy") // 不影响原对象（因为 copy 是值）
	fmt.Println("copy:", copy, "orig:", *u)
} 
