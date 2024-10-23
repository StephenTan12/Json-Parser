package main

import (
	"fmt"
	"io"
	"os"
)

type JsonTokenizer struct {
	filePtr         *os.File
	currFilePointer int64
	buffer          []byte
	tokens          []string
}

func (j *JsonTokenizer) Init(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	j.filePtr = file
	j.currFilePointer = 0
	j.buffer = make([]byte, 1)
	j.tokens = make([]string, 0)
	return nil
}

func (j *JsonTokenizer) Close() {
	j.filePtr.Close()
}

func (j *JsonTokenizer) BuildTokens() {
	_, err := j.filePtr.ReadAt(j.buffer, j.currFilePointer)
	if err != nil {
		if err == io.EOF {
			j.tokens = append(j.tokens, "\b255")
			return
		}
		fmt.Println(err)
		return
	}

	char := rune(j.buffer[0])
	j.currFilePointer += 1

	switch char {
	case LBraces:
		j.tokens = append(j.tokens, string(LBraces))
		j.BuildTokens()
	case RBraces:
		j.tokens = append(j.tokens, string(RBraces))
		j.BuildTokens()
	case LBracket:
		j.tokens = append(j.tokens, string(LBracket))
		j.BuildTokens()
	case RBracket:
		j.tokens = append(j.tokens, string(RBracket))
		j.BuildTokens()
	case Quotation:
		value, _ := j.GetValue(true)
		if len(value) != 0 {
			j.tokens = append(j.tokens, string(Quotation), value, string(Quotation))
		} else {
			j.tokens = append(j.tokens, string(Quotation), string(Quotation))
		}
		j.BuildTokens()
	case Colon:
		j.tokens = append(j.tokens, string(Colon))
		j.BuildTokens()
	case Comma:
		j.tokens = append(j.tokens, string(Comma))
		j.BuildTokens()
	case Space:
		j.BuildTokens()
	case Tab:
		j.BuildTokens()
	case Newline:
		j.BuildTokens()
	default:
		j.currFilePointer -= 1
		value, _ := j.GetValue(false)
		j.tokens = append(j.tokens, value)
		j.BuildTokens()
	}
}

func (j *JsonTokenizer) GetValue(fromQuotes bool) (string, string) {
	value := make([]byte, 0)

	for {
		_, err := j.filePtr.ReadAt(j.buffer, j.currFilePointer)
		if err != nil {
			if err == io.EOF {
				return string(value), ""
			}
			fmt.Println(err)
			return string(value), ""
		}

		char := rune(j.buffer[0])

		if fromQuotes && char == '"' {
			j.currFilePointer += 1
			return string(value), string(Quotation)
		}

		if !fromQuotes && (j.IsObject(char) || j.IsPunctuation(char) || j.IsWhitespace(char)) {
			return string(value), ""
		}

		j.currFilePointer += 1
		value = append(value, j.buffer[0])
	}
}

func (j *JsonTokenizer) IsWhitespace(char rune) bool {
	return char == Space || char == Tab || char == Newline
}

func (j *JsonTokenizer) IsPunctuation(char rune) bool {
	return char == Quotation || char == Comma || char == Colon
}

func (j *JsonTokenizer) IsObject(char rune) bool {
	return char == LBraces || char == RBraces || char == LBracket || char == RBracket
}
