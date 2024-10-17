package main

import "fmt"

func main() {
	j := JsonTokenizer{}
	j.Init("./app/test.json")
	j.BuildTokens()

	fmt.Println(j.tokens)

	parser := JsonParser{
		tokens:           &j.tokens,
		currTokenPointer: 0,
	}
	isValid := parser.ParseTokens()

	fmt.Println(isValid)

	j.Close()
}
