package main

import "fmt"

const englishPrefix = "Hello "
const chinesePrefix = "你好 "
const spanishPrefix = "Hola "
const frenchPrefix = "Bonjour "

func Hello(name string, language string) string {
	if name == "" {
		name = "World"
	}

	return greetingPrefix(language) + name
}

func greetingPrefix(language string) string {
	switch language {
	case "Chinese":
		return chinesePrefix
	case "Spanish":
		return spanishPrefix
	case "French":
		return frenchPrefix
	default:
		return englishPrefix
	}
}

func main() {
	fmt.Println(Hello("Fan", ""))
}
