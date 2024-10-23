package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("a path to a file needs to be provided")
		return
	} else if len(args) > 2 {
		fmt.Println("too many arguments provided")
		return
	}
	j := JsonTokenizer{}
	err := j.Init(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
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
	} else {
		fmt.Println("invalid json")
	}

	j.Close()
}
