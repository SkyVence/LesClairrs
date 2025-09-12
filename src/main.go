package main

import (
	"fmt"
)

func main() {
	lang, _ := Load("fr")
	title := lang.Text("menu.title")
	welcome := lang.Text("menu.welcome")
	fmt.Print(title)
	fmt.Print(welcome)
}
 