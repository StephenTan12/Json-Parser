package main

import "fmt"

func main() {
	j := JsonTokenizer{}
	j.Init("./app/test.json")
	j.BuildTokens()
	fmt.Println(j.tokens)
	j.Close()
}
