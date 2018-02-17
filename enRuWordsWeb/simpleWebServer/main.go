package main

import (
	"bufio"
	"os"
)

func main() {
	println("куку")
	reader := bufio.NewReader(os.Stdin)
	for {
		println("Как тебя зовут")
		text, _ := reader.ReadString('\n')
		println("Привет,", text)
	}
}
