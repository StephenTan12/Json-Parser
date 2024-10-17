package main

import (
	"fmt"
	"io"
)

type JsonParser struct {
	tokens           *[]string
	currTokenPointer int
	currStack        []rune
	currStackPointer int
}

type Error struct {
	s string
}

func (e *Error) Error() string {
	return e.s
}

func (j *JsonParser) ParseTokens() bool {
	j.currStack = make([]rune, len(*j.tokens))
	j.currStackPointer = 0

	err := j.parseEntry()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if j.currStackPointer == 0 && j.currStack[0] == 0 && j.currTokenPointer == len(*j.tokens) {
		return true
	}

	return false
	/*

		LBraces
		-> RBraces | LBracket | Quotes

		LBracket
		-> RBracket | LBraces | Quotes | value

		RBracket
		-> RBracket, RBraces, Comma

		RBraces
		-> RBracket, RBraces, Comma

		Quotes Value Quotes
		-> RBraces, RBracket, Comma, Colon

		Comma:
		-> LBraces, LBracket, Quotes, Value

		Colon:
		-> LBraces, LBracket, Quotes, Value

		Value:
		-> RBraces, RBracket, Comma
	*/
}

func (j *JsonParser) parseEntry() error {
	token, err := j.nextToken()
	if err != nil {
		return err
	}

	switch rune(token[0]) {
	case LBraces:
		j.currStack[j.currStackPointer] = LBraces
		j.currStackPointer += 1

		err := j.parseLBraces()
		if err != nil {
			return err
		}
	case LBracket:
		j.currStack[j.currStackPointer] = LBracket
		j.currStackPointer += 1

		err := j.parseLBracket()
		if err != nil {
			return err
		}
	default:
		return &Error{s: "invalid json"}
	}

	return nil
}

func (j *JsonParser) parseLBraces() error {
	token, err := j.nextToken()
	if err != nil {
		return err
	}

	switch rune(token[0]) {
	case RBraces:
		if j.currStack[j.currStackPointer-1] != LBraces {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBraces()
		if err != nil {
			return err
		}
	case LBracket:
		j.currStack[j.currStackPointer] = LBracket
		j.currStackPointer += 1

		err := j.parseLBracket()
		if err != nil {
			return err
		}
	case Quotation:
		token, err = j.nextToken()
		if err != nil {
			return err
		}

		if len(token) > 1 || (rune(token[0]) != Quotation && len(token) == 1) {
			fmt.Println("Got value: " + string(token))
			token, err = j.nextToken()
			if err != nil {
				return err
			}
		}

		if !(len(token) == 1 && rune(token[0]) == Quotation) {
			return &Error{s: "invalid json"}
		}

		err := j.parseQuotations()
		if err != nil {
			return err
		}
	default:
		return &Error{s: "invalid json"}
	}

	return nil
}

func (j *JsonParser) parseLBracket() error {
	token, err := j.nextToken()
	if err != nil {
		return err
	}

	switch rune(token[0]) {
	case RBracket:
		if j.currStack[j.currStackPointer-1] != LBracket {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBracket()
		if err != nil {
			return err
		}
	case LBraces:
		j.currStack[j.currStackPointer] = LBraces
		j.currStackPointer += 1

		err := j.parseLBraces()
		if err != nil {
			return err
		}
	case Quotation:
		token, err = j.nextToken()
		if err != nil {
			return err
		}

		if len(token) > 1 || (rune(token[0]) != Quotation && len(token) == 1) {
			fmt.Println("Got value: " + string(token))
			token, err = j.nextToken()
			if err != nil {
				return err
			}
		}

		if !(len(token) == 1 && rune(token[0]) == Quotation) {
			return &Error{s: "invalid json"}
		}

		err := j.parseQuotations()
		if err != nil {
			return err
		}
	default:
		fmt.Println("Got value: " + token)
		err := j.parseValue()
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JsonParser) parseRBracket() error {
	token, err := j.nextToken()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	switch rune(token[0]) {
	case RBracket:
		if j.currStack[j.currStackPointer-1] != LBracket {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBracket()
		if err != nil {
			return err
		}
	case RBraces:
		if j.currStack[j.currStackPointer-1] != LBraces {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBraces()
		if err != nil {
			return err
		}
	case Comma:
		fmt.Println("comma")
		err := j.parseComma()
		if err != nil {
			return err
		}
	default:
		fmt.Println("Got value: " + string(token))
		err := j.parseValue()
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JsonParser) parseRBraces() error {
	token, err := j.nextToken()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	switch rune(token[0]) {
	case RBracket:
		if j.currStack[j.currStackPointer-1] != LBracket {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBracket()
		if err != nil {
			return err
		}
	case RBraces:
		if j.currStack[j.currStackPointer-1] != LBraces {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBraces()
		if err != nil {
			return err
		}
	case Comma:
		fmt.Println("comma")
		err := j.parseComma()
		if err != nil {
			return err
		}
	default:
		return &Error{s: "invalid json"}
	}

	return nil
}

func (j *JsonParser) parseComma() error {
	token, err := j.nextToken()
	if err != nil {
		return err
	}

	switch rune(token[0]) {
	case LBraces:
		j.currStack[j.currStackPointer] = LBraces
		j.currStackPointer += 1

		err := j.parseLBraces()
		if err != nil {
			return err
		}
	case LBracket:
		j.currStack[j.currStackPointer] = LBracket
		j.currStackPointer += 1

		err := j.parseLBracket()
		if err != nil {
			return err
		}
	case Quotation:
		token, err = j.nextToken()
		if err != nil {
			return err
		}

		if len(token) > 1 || (rune(token[0]) != Quotation && len(token) == 1) {
			fmt.Println("Got value: " + string(token))
			token, err = j.nextToken()
			if err != nil {
				return err
			}
		}

		if !(len(token) == 1 && rune(token[0]) == Quotation) {
			return &Error{s: "invalid json"}
		}

		err := j.parseQuotations()
		if err != nil {
			return err
		}
	default:
		fmt.Println("Got value: " + string(token))
		err := j.parseValue()
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JsonParser) parseColon() error {
	token, err := j.nextToken()
	if err != nil {
		return err
	}

	switch rune(token[0]) {
	case LBraces:
		j.currStack[j.currStackPointer] = LBraces
		j.currStackPointer += 1

		err := j.parseLBraces()
		if err != nil {
			return err
		}
	case LBracket:
		j.currStack[j.currStackPointer] = LBracket
		j.currStackPointer += 1

		err := j.parseLBracket()
		if err != nil {
			return err
		}
	case Quotation:
		token, err = j.nextToken()
		if err != nil {
			return err
		}

		if len(token) > 1 || (rune(token[0]) != Quotation && len(token) == 1) {
			fmt.Println("Got value: " + string(token))
			token, err = j.nextToken()
			if err != nil {
				return err
			}
		}

		if !(len(token) == 1 && rune(token[0]) == Quotation) {
			return &Error{s: "invalid json"}
		}

		err := j.parseQuotations()
		if err != nil {
			return err
		}
	default:
		fmt.Println("Got value: " + string(token))
		err := j.parseValue()
		if err != nil {
			return err
		}
	}

	return nil
}

func (j *JsonParser) parseQuotations() error {
	token, err := j.nextToken()
	if err != nil {
		return err
	}

	switch rune(token[0]) {
	case RBracket:
		if j.currStack[j.currStackPointer-1] != LBracket {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBracket()
		if err != nil {
			return err
		}
	case RBraces:
		if j.currStack[j.currStackPointer-1] != LBraces {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBraces()
		if err != nil {
			return err
		}
	case Comma:
		fmt.Println("comma")
		err := j.parseComma()
		if err != nil {
			return err
		}
	case Colon:
		fmt.Println("colon  ", token)
		err := j.parseColon()
		if err != nil {
			return err
		}
	default:
		return &Error{s: "invalid json"}
	}

	return nil
}

func (j *JsonParser) parseValue() error {
	token, err := j.nextToken()
	if err != nil {
		return err
	}

	switch rune(token[0]) {
	case RBracket:
		if j.currStack[j.currStackPointer-1] != LBracket {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBracket()
		if err != nil {
			return err
		}
	case RBraces:
		if j.currStack[j.currStackPointer-1] != LBraces {
			return &Error{s: "invalid json"}
		}
		j.currStack[j.currStackPointer-1] = 0
		j.currStackPointer -= 1

		err := j.parseRBraces()
		if err != nil {
			return err
		}
	case Comma:
		fmt.Println("comma")
		err := j.parseComma()
		if err != nil {
			return err
		}
	default:
		return &Error{s: "invalid json"}
	}

	return nil
}

func (j *JsonParser) nextToken() (string, error) {
	if j.currTokenPointer >= len(*j.tokens) {
		return "", io.EOF
	}
	token := (*j.tokens)[j.currTokenPointer]
	j.currTokenPointer += 1
	fmt.Println(token)

	return token, nil
}
