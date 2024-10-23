package main

import (
	"fmt"
	"io"
	"strconv"
)

type JsonParser struct {
	tokens           []string
	currTokenPointer int
	currStack        []rune
	currStackPointer int
}

func (e *Error) Error() string {
	return e.s
}

func (j *JsonParser) ParseTokens() (Json, bool) {
	j.currStack = make([]rune, len(j.tokens))
	j.currStackPointer = 0

	object, err := j.parseEntry()
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	if j.currStackPointer == 0 && j.currStack[0] == 0 && j.currTokenPointer == len(j.tokens)-1 {
		return object, true
	}
	return nil, false
}

func (j *JsonParser) parseEntry() (Json, error) {
	token, err := j.nextToken()
	if err != nil {
		return nil, err
	}

	switch rune(token[0]) {
	case LBraces:
		j.currStack[j.currStackPointer] = LBraces
		j.currStackPointer += 1

		obj, err := j.parseLBraces()
		if err != nil {
			return nil, err
		}
		return obj, err
	case LBracket:
		j.currStack[j.currStackPointer] = LBracket
		j.currStackPointer += 1

		obj, err := j.parseLBracket()
		if err != nil {
			return nil, err
		}
		return obj, err
	default:
		return nil, &Error{s: "invalid json"}
	}
}

func (j *JsonParser) parseLBraces() (JsonObject, error) {
	object := JsonObject{}
	already_seen_comma := false
out:
	for {
		token, err := j.nextToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		char := rune(token[0])

		switch char {
		case RBraces:
			if j.removeFromStack(LBraces) != nil {
				return nil, err
			}
			break out
		case Quotation:
			already_seen_comma = false

			key, value, err := j.parseKeyValuePair()
			if err != nil {
				return nil, err
			}
			object[key] = value

			// check for future values?
			token, err = j.nextToken()
			if err != nil {
				return nil, err
			}
			char = rune(token[0])

			switch char {
			case Comma:
				if already_seen_comma {
					return nil, &Error{s: "invalid json"}
				}
				already_seen_comma = true
				continue
			case RBraces:
				if j.currStack[j.currStackPointer-1] != LBraces {
					return nil, &Error{s: "invalid json"}
				}
				j.currStack[j.currStackPointer-1] = 0
				j.currStackPointer -= 1
				break out
			default:
				return nil, &Error{s: "invalid json"}
			}
		default:
			return nil, &Error{s: "invalid json"}
		}
	}

	if already_seen_comma {
		return nil, &Error{s: "invalid json"}
	}

	return object, nil
}

func (j *JsonParser) parseLBracket() (JsonArray, error) {
	object := JsonArray{}
	already_seen_comma := false
out:
	for {
		token, err := j.nextToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(token) > 1 {
			if token == "true" {
				object = append(object, true)
			} else if token == "false" {
				object = append(object, false)
			} else {
				value, err := strconv.Atoi(token)
				if err != nil {
					return nil, err
				}
				object = append(object, value)
			}
			already_seen_comma = false
			continue
		}

		char := rune(token[0])
		switch char {
		case RBracket:
			if j.removeFromStack(LBracket) != nil {
				return nil, err
			}
			break out
		case Quotation:
			// getting element
			element, err := j.parseString()
			if err != nil {
				return nil, err
			}
			object = append(object, element)
		case Comma:
			if len(object) < 1 || already_seen_comma {
				return nil, &Error{s: "invalid json"}
			}
			already_seen_comma = true
			continue
		case LBracket:
			j.currStack[j.currStackPointer] = LBracket
			j.currStackPointer += 1
			newObject, err := j.parseLBracket()
			if err != nil {
				return nil, err
			}
			object = append(object, newObject)
		case LBraces:
			j.currStack[j.currStackPointer] = LBraces
			j.currStackPointer += 1
			newObject, err := j.parseLBraces()
			if err != nil {
				return nil, err
			}
			object = append(object, newObject)
		default:
			return nil, &Error{s: "invalid json"}
		}
		already_seen_comma = false
	}

	if already_seen_comma {
		return nil, &Error{s: "invalid json"}
	}

	return object, nil
}

func (j *JsonParser) parseKeyValuePair() (string, Json, error) {
	// getting key
	key, err := j.parseString()
	if err != nil {
		return "", nil, err
	}

	token, err := j.nextToken()
	if err != nil {
		return "", nil, err
	}
	if rune(token[0]) != Colon {
		return "", nil, err
	}

	// getting value
	value, err := j.parseValue()
	if err != nil {
		return "", nil, err
	}
	return key, value, err
}

func (j *JsonParser) parseString() (string, error) {
	token, err := j.nextToken()
	if err != nil {
		return "", err
	}

	key := ""
	if len(token) > 1 || (rune(token[0]) != Quotation && len(token) == 1) {
		key = string(token)
		token, err = j.nextToken()
		if err != nil {
			return "", err
		}
	}
	if !(len(token) == 1 && rune(token[0]) == Quotation) {
		return "", &Error{s: "invalid json"}
	}

	return key, nil
}

func (j *JsonParser) parseValue() (Json, error) {
	token, err := j.nextToken()
	if err != nil {
		return "", err
	}

	if len(token) > 1 {
		if token == "true" {
			return true, nil
		} else if token == "false" {
			return false, nil
		} else {
			value, err := strconv.Atoi(token)
			if err != nil {
				return "", err
			}
			return value, nil
		}
	}

	char := rune(token[0])
	switch char {
	case LBraces:
		j.currStack[j.currStackPointer] = LBraces
		j.currStackPointer += 1
		value, err := j.parseLBraces()
		if err != nil {
			return nil, err
		}
		return value, nil
	case LBracket:
		j.currStack[j.currStackPointer] = LBracket
		j.currStackPointer += 1

		newObj, err := j.parseLBracket()
		if err != nil {
			return nil, err
		}
		return newObj, nil
	case Quotation:
		token, err = j.nextToken()
		if err != nil {
			return nil, err
		}
		value := ""
		if len(token) > 1 || (rune(token[0]) != Quotation && len(token) == 1) {
			value = string(token)
			token, err = j.nextToken()
			if err != nil {
				return nil, err
			}
		}
		if !(len(token) == 1 && rune(token[0]) == Quotation) {
			return nil, &Error{s: "invalid json"}
		}
		return value, err
	default:
		return nil, &Error{s: "invalid json"}
	}
}

func (j *JsonParser) removeFromStack(element rune) error {
	if j.currStack[j.currStackPointer-1] != element {
		return &Error{s: "invalid json"}
	}
	j.currStack[j.currStackPointer-1] = 0
	j.currStackPointer -= 1
	return nil
}

func (j *JsonParser) nextToken() (string, error) {
	token := (j.tokens)[j.currTokenPointer]
	fmt.Println("token: ", token)
	if token == "\b255" {
		j.currTokenPointer += 1
		return "", io.EOF
	}
	j.currTokenPointer += 1
	return token, nil
}
