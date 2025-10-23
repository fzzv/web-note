package main

import (
	"log"
	"net/http"
)

type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func (i *InMemoryPlayerStore) RecordWin(name string) {}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}}
	handler := http.HandlerFunc(server.ServeHTTP)
	if err := http.ListenAndServe(":5000", handler); err != nil {
		log.Fatal(err)
	}
}
