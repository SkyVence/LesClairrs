package main

import (
	"fmt"
)

func main() {
	lang, _ := Load("fr")
	title := lang.Text("menu.title")
	fmt.Print(title)
}
