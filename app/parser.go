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

	for {
		token, err := j.nextToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		char := rune(token[0])

		if char == RBraces {
			if j.currStack[j.currStackPointer-1] != LBraces {
				return nil, &Error{s: "invalid json"}
			}
			j.currStack[j.currStackPointer-1] = 0
			j.currStackPointer -= 1
			break
		} else if char == Quotation {
			already_seen_comma = false
			// getting key
			token, err = j.nextToken()
			if err != nil {
				return nil, err
			}

			key := ""
			if len(token) > 1 || (rune(token[0]) != Quotation && len(token) == 1) {
				key = string(token)
				token, err = j.nextToken()
				if err != nil {
					return nil, err
				}
			}
			if !(len(token) == 1 && rune(token[0]) == Quotation) {
				return nil, &Error{s: "invalid json"}
			}

			token, err = j.nextToken()
			if err != nil {
				return nil, err
			}
			if rune(token[0]) != Colon {
				return nil, err
			}

			token, err = j.nextToken()
			if err != nil {
				return nil, err
			}

			// getting value
			char = rune(token[0])
			if len(token) > 1 {
				if token == "true" {
					object[key] = true
				} else if token == "false" {
					object[key] = false
				} else {
					value, err := strconv.Atoi(token)
					if err != nil {
						return nil, err
					}
					object[key] = value
				}
			} else if char == LBraces {
				j.currStack[j.currStackPointer] = LBraces
				j.currStackPointer += 1
				value, err := j.parseLBraces()
				if err != nil {
					return nil, err
				}
				object[key] = value
			} else if char == LBracket {
				j.currStack[j.currStackPointer] = LBracket
				j.currStackPointer += 1

				newObj, err := j.parseLBracket()
				if err != nil {
					return nil, err
				}
				object[key] = newObj
			} else if char == Quotation {
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
				object[key] = value
			} else {
				return nil, &Error{s: "invalid json"}
			}

			// check for future values?
			token, err = j.nextToken()
			if err != nil {
				return nil, err
			}
			char = rune(token[0])
			if char == Comma {
				if already_seen_comma {
					return nil, &Error{s: "invalid json"}
				}
				already_seen_comma = true
				continue
			} else if char == RBraces {
				if j.currStack[j.currStackPointer-1] != LBraces {
					return nil, &Error{s: "invalid json"}
				}
				j.currStack[j.currStackPointer-1] = 0
				j.currStackPointer -= 1
				break
			} else {
				return nil, &Error{s: "invalid json"}
			}
		} else {
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

	for {
		token, err := j.nextToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		char := rune(token[0])

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
		} else if char == RBracket {
			if j.currStack[j.currStackPointer-1] != LBracket {
				return nil, &Error{s: "invalid json"}
			}
			j.currStack[j.currStackPointer-1] = 0
			j.currStackPointer -= 1
			break
		} else if char == Quotation {
			// getting element
			token, err = j.nextToken()
			if err != nil {
				return nil, err
			}
			element := ""
			if len(token) > 1 || (rune(token[0]) != Quotation && len(token) == 1) {
				element = string(token)
				token, err = j.nextToken()
				if err != nil {
					return nil, err
				}
			}
			if !(len(token) == 1 && rune(token[0]) == Quotation) {
				return nil, &Error{s: "invalid json"}
			}
			object = append(object, element)
		} else if char == Comma {
			if len(object) < 1 || already_seen_comma {
				return nil, &Error{s: "invalid json"}
			}
			already_seen_comma = true
			continue
		} else if char == LBracket {
			j.currStack[j.currStackPointer] = LBracket
			j.currStackPointer += 1
			newObject, err := j.parseLBracket()
			if err != nil {
				return nil, err
			}
			object = append(object, newObject)
		} else if char == LBraces {
			j.currStack[j.currStackPointer] = LBraces
			j.currStackPointer += 1
			newObject, err := j.parseLBraces()
			if err != nil {
				return nil, err
			}
			object = append(object, newObject)
		} else {
			return nil, &Error{s: "invalid json"}
		}
		already_seen_comma = false
	}

	if already_seen_comma {
		return nil, &Error{s: "invalid json"}
	}

	return object, nil
}

func (j *JsonParser) nextToken() (string, error) {
	token := (j.tokens)[j.currTokenPointer]
	if token == "\b255" {
		j.currTokenPointer += 1
		return "", io.EOF
	}
	j.currTokenPointer += 1
	return token, nil
}
