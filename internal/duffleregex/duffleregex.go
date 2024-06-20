package duffleregex

import (
	"fmt"
	"log"

	"github.com/tflexsoom/duffle/internal/container"
)

const (
	D_REGEX_NONE    uint8 = 0
	D_REGEX_SUCCESS uint8 = 1
	D_REGEX_FAIL    uint8 = 2
	D_REGEX_CLASS   uint8 = 4
)

type DuffleRegex struct {
	regex string
	index int
	class []rune
	flag  uint8
}

func NewDuffleWithAssert(regex string) (DuffleRegex, error) {
	parenthesis := 0
	squarebracket := false
	negative := false
	hasLast := false
	escaped := false

	escapable := []rune{'\\', '(', ')', '[', ']', '*', '+', '.'}

	for i, v := range regex {
		if escaped && container.In(v, escapable) {
			hasLast = true
			escaped = false
			continue
		} else if escaped {
			return DuffleRegex{}, fmt.Errorf("unknown escape character %v", v)
		}

		switch v {
		case '\\':
			escaped = true
		case '(':
			parenthesis += 1
		case ')':
			parenthesis -= 1
			if parenthesis < 0 {
				return DuffleRegex{}, fmt.Errorf("unmatched parenthesis at position %v", i)
			}
			hasLast = true
		case '*':
			if !hasLast {
				return DuffleRegex{}, fmt.Errorf("postfix operator has nothing to postfix at %v", i)
			}
			hasLast = false
		case '+':
			if !hasLast {
				return DuffleRegex{}, fmt.Errorf("postfix operator has nothing to postfix at %v", i)
			}
			hasLast = false
		case '[':
			if squarebracket {
				return DuffleRegex{}, fmt.Errorf("unexpected '[' at position %v", i)
			}
			squarebracket = true
		case ']':
			if !squarebracket {
				return DuffleRegex{}, fmt.Errorf("unexpected ']' at position %v", i)
			}
			squarebracket = false
			negative = false
		case '^':
			if squarebracket && negative {
				return DuffleRegex{}, fmt.Errorf("too many '^' at position %v", i)
			}
			negative = true
		case '.':
			hasLast = true
		default:
			hasLast = true
		}
	}

	return NewDuffleRegex(regex), nil
}

func NewDuffleRegex(regex string) DuffleRegex {
	return DuffleRegex{
		regex: regex,
		index: 0,
		class: []rune{},
		flag:  D_REGEX_NONE,
	}
}

func (dr DuffleRegex) IsSuccess() bool {
	return dr.flag&D_REGEX_SUCCESS != 0
}

func (dr DuffleRegex) IsFail() bool {
	return dr.flag&D_REGEX_FAIL != 0
}

func (dr DuffleRegex) Step(character rune) ([]DuffleRegex, error) {
	return []DuffleRegex{}, nil
}

type matchStackTrace struct {
	regexes    []DuffleRegex
	inputIndex int
}

type regexOutput uint8

const (
	REGEX_OUTPUT_NONE regexOutput = iota
	REGEX_OUTPUT_SUCCESS
	REGEX_OUTPUT_FAIL
)

func (dr DuffleRegex) subMatch(input string) regexOutput {
	length := len(input)
	stack := make([]matchStackTrace, 0, 16)

	for i, v := range input {
		regexes, err := dr.Step(v)

		if err != nil {
			log.Printf("error %v found at input position %d", err, i)
			return REGEX_OUTPUT_FAIL
		}

		if i < length-1 {
			stack = append(stack, matchStackTrace{
				regexes:    regexes,
				inputIndex: i,
			})
		}

		if dr.IsFail() {
			return REGEX_OUTPUT_FAIL
		}
	}

	if dr.IsSuccess() {
		return REGEX_OUTPUT_SUCCESS
	}

	stackLength := len(stack)
	for i := stackLength - 1; i >= 0; i-- {
		regexLength := len(stack[i].regexes)
		inputSlice := stack[i].inputIndex
		for j := 0; j < regexLength; j++ {
			output := stack[i].regexes[j].subMatch(input[inputSlice:length])
			if output != REGEX_OUTPUT_NONE {
				return output
			}
		}
	}

	return REGEX_OUTPUT_NONE
}

func (dr DuffleRegex) Match(input string) bool {
	return dr.subMatch(input) == REGEX_OUTPUT_SUCCESS
}
