package main

import "fmt"

func main() {
	j := JsonTokenizer{}
	j.Init("./app/test.json")
	j.BuildTokens()

	fmt.Println("total tokens:", len(j.tokens))
	fmt.Println("tokens: ", j.tokens)

	parser := JsonParser{
		tokens:           j.tokens,
		currTokenPointer: 0,
	}
	object, isValid := parser.ParseTokens()
	if isValid {
		fmt.Println("valid json: ", object)
	}

	j.Close()
}
