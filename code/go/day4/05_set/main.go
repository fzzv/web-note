package main

import "fmt"

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] { return make(Set[T]) }
func (s Set[T]) Add(v T)          { s[v] = struct{}{} }
func (s Set[T]) Has(v T) bool     { _, ok := s[v]; return ok }
func (s Set[T]) Union(other Set[T]) Set[T] {
	out := NewSet[T]()
	for v := range s { out.Add(v) }
	for v := range other { out.Add(v) }
	return out
}

func main() {
	a := NewSet[int](); a.Add(1); a.Add(2)
	b := NewSet[int](); b.Add(2); b.Add(3)
	u := a.Union(b)
	fmt.Println("has 1:", u.Has(1), "has 3:", u.Has(3))
} 
