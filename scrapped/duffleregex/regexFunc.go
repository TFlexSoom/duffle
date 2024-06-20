package duffleregex

import (
	"fmt"
	"log"

	"github.com/tflexsoom/duffle/internal/container"
)

func NewDuffleWithAssert(regex string) (DuffleRegex, error) {
	parenthesis := 0
	squarebracket := false
	curly := false
	negative := false
	hasLast := false
	escaped := false
	escapable := []rune{'\\', '(', ')', '[', ']', '{', '}', '*', '+', '.', '?'}
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
		case '?':
			if !hasLast {
				return DuffleRegex{}, fmt.Errorf("postfix operator has nothing to postfix at %v", i)
			}
			hasLast = false
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
		case '{':
			// TODO add check logic to parsing logic
			if curly {
				fmt.Errorf("Unmatched '}' at position %v", i)
			}
			curly = true
		case '}':
			if !curly {
				fmt.Errorf("Unexpected '}' at position %v", i)
			}
		case '.':
			hasLast = true
		default:
			hasLast = true
		}
	}
	return NewDuffleRegex(regex), nil
}
func NewDuffleRegex(regex string) DuffleRegex {
	dr := DuffleRegex{
		input:  []rune{},
		regex:  []rune(regex),
		states: make([]state, 0, len(regex)>>2),
		rules:  make([]rule, 0, len(regex)),
		groups: make(map[uint]uint, len(regex)),
		trails: make([]uint, 1, 4),
	}

	dr.exploreState(dr.trails[0])

	return dr
}

func (dr *DuffleRegex) MatchRunify(input string) bool {
	return dr.Match([]rune(input))
}

func (dr *DuffleRegex) Match(input []rune) bool {
	result := REGEX_RESULT_NONE
	dr.input = input
	result, err := dr.Resolve()
	if err != nil {
		log.Printf("error from resolving regex: %v", err)
	}

	return result == REGEX_RESULT_SUCCESS
}

func (dr *DuffleRegex) Step(character rune) (RegexResult, error) {
	dr.input = append(dr.input, character)
	return dr.Resolve()
}

func (dr *DuffleRegex) Resolve() (RegexResult, error) {
	currentState := dr.states[dr.current]

	if len(currentState.transitions) <= 0 {
		return currentState.result, nil
	}

	rule := currentState.transitions[0]

	// 2 If single char
	//   eq move forward index
	//   neq fail
	//  If wild card
	//   move forward index
	//  If Class check contains char
	//   contains forward index
	//   disjoint fail
	//  If Neg class opposite of 3
	//    ...
	//  If x-y check above or next rule
	//  If zero-one ...
	//  If zero-more ...
	//  IF one-more ...
	//  If Union check between 2 rules until one fails
	//
	//  If group check look for other group and make current slice

	return REGEX_RESULT_NONE, nil
}

func (dr *DuffleRegex) getCharClass(regexIndex uint) ([]rune, bool) {
	regex := dr.regex
	regexlen := uint(len(regex))
	contains := make([]rune, 0, 16)
	flag := false
	for i := regexIndex + 1; i < regexlen; i++ {
		if regex[i] == '\\' {
			flag = true
			continue
		} else if regex[i] == ']' && !flag {
			break
		}

		flag = false
		contains = append(contains, regex[i])
	}

	if contains[0] == '^' {
		return contains[1:], true
	}

	return contains, false
}

func (dr *DuffleRegex) gotoStateWithIndex(regexIndex uint, inputIndex uint) {
	rules, err := exploreRules(regexIndex)
	if err != nil {
		log.Printf("error occurred in searching state %v", err)
		return
	}

	dr.states = append(dr.states, state{
		inputIndex:  inputIndex,
		result:      REGEX_RESULT_NONE,
		transitions: rules,
	})
}

func (dr *DuffleRegex) exploreRules(regexIndex uint) ([]rule, error) {
	if dr.rules[regexIndex] != nil {
		return {[]dr.rules{}// already explored
	}

	regexLength := uint(len(dr.regex))
	if regexIndex >= regexLength {
		log.Printf("bad regexIndex past length")
		return // Bad param
	}

	rules := dr.exploreRules(regexIndex)
	firstRune := dr.regex[regexIndex]
	var contains []rune
	flags := RULE_FLAG_NONE
	switch firstRune {
	case '\\':
		if regexIndex+1 > regexLength {
			log.Printf("bad '\\' character in regex at pos %v", regexIndex)
			return
		}
		contains = []rune{dr.regex[regexIndex+1]}
		regexIndex += 2
	case '[':
		var isNeg bool
		contains, isNeg = dr.getCharClass(regexIndex)
		regexIndex += uint(len(contains)) + 2
		if isNeg {
			flags |= RULE_FLAG_NEG
			regexIndex += 1
		}
	case '(':
		contains, firstIndex, altIndex := dr.getSubRegex(regexIndex)
		regexIndex = altIndex
	case '.':
		contains = []rune{}
		flags |= RULE_FLAG_NEG
		regexIndex += 1
	default:
		contains = []rune{firstRune}
		regexIndex += 1
	}
}
