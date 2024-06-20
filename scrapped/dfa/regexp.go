package dfa

// import (
// 	"errors"

// 	"github.com/tflexsoom/duffle/internal/container"
// )

// type ruleFlag int8

// const (
// 	RULE_SINGLE_RUNE ruleFlag = ruleFlag(1)
// 	RULE_GROUP       ruleFlag = ruleFlag(2)
// 	RULE_UNION       ruleFlag = ruleFlag(4)
// 	RULE_INF_UPPER   ruleFlag = ruleFlag(8)
// 	RULE_CHAR_CLASS  ruleFlag = ruleFlag(16)
// 	RULE_NEG         ruleFlag = ruleFlag(32)
// )

// type rule struct {
// 	val        rune
// 	lowerBound uint
// 	upperBound uint
// 	flags      ruleFlag
// }

// func (r rule) isNoneRule() bool {
// 	return r.flags == 0
// }

// func (r rule) isSingleRune() bool {
// 	return r.flags&RULE_SINGLE_RUNE != 0
// }

// func (r rule) isGroupRune() bool {
// 	return r.flags&RULE_GROUP != 0
// }

// func (r rule) isUnion() bool {
// 	return r.flags&RULE_UNION != 0
// }

// func (r rule) isInfiniteUpperRange() bool {
// 	return r.flags&RULE_INF_UPPER != 0
// }

// func (r rule) isCharacterClass() bool {
// 	return r.flags&RULE_CHAR_CLASS != 0
// }

// func (r rule) isNegative() bool {
// 	return r.flags&RULE_NEG != 0
// }

// func NewRuneRule(val rune) rule {
// 	return rule{
// 		val:        val,
// 		lowerBound: 0,
// 		upperBound: 0,
// 		flags:      RULE_SINGLE_RUNE,
// 	}
// }

// func NewSimpleGroupRule() rule {
// 	return NewRangedGroupRule(0, 0)
// }

// func NewRangedGroupRule(lower uint, upper uint) rule {
// 	return NewGroupRule(lower, upper, false)
// }

// func NewNegativeGroupRule() rule {
// 	return NewGroupRule(0, 0, true)
// }

// func NewGroupRule(lower uint, upper uint, isNeg bool) rule {
// 	flags := RULE_GROUP
// 	if isNeg {
// 		flags |= RULE_NEG
// 	}

// 	return rule{
// 		val:        rune(0),
// 		lowerBound: lower,
// 		upperBound: upper,
// 		flags:      flags,
// 	}

// }

// func NewStarGroupRule() rule {
// 	return rule{
// 		val:        rune(0),
// 		lowerBound: 0,
// 		upperBound: 0,
// 		flags:      RULE_GROUP | RULE_INF_UPPER,
// 	}
// }

// func NewPlusGroupRule() rule {
// 	return rule{
// 		val:        rune(0),
// 		lowerBound: 1,
// 		upperBound: 0,
// 		flags:      RULE_GROUP | RULE_INF_UPPER,
// 	}
// }

// func NewUnionRule() rule {
// 	return rule{
// 		val:        rune(0),
// 		lowerBound: 0,
// 		upperBound: 0,
// 		flags:      RULE_GROUP | RULE_UNION,
// 	}
// }

// type RegexpTree container.Tree[rule]

// type slice struct {
// 	from uint
// 	to   uint
// }

// func FromRegexp(expression string) (RegexpTree, error) {
// 	slices, err := getGroups(expression)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return fromRegexpAndGroups(expression, append(slices, slice{from: 0, to: uint(len(expression))}))
// }

// func fromRegexpAndGroups(expression string, groups []slice) (RegexpTree, error) {
// 	result := container.NewGraphTreeCap[rule](8, uint(len(expression)))
// 	result.SetValue(NewSimpleGroupRule())

// 	temp := rule{}

// 	addRule := func(r rule) {
// 		if !temp.isNoneRule() {
// 			result.AddChild(temp)
// 		}

// 		temp = r
// 	}

// 	group_index := len(groups) - 1
// 	group := groups[group_index]

// 	range_flag := false
// 	starting_range := uint(0)
// 	bracket_flag := false
// 	starting_bracket := uint(0)

// 	for i, v := range expression {
// 		switch v {
// 		case '|':
// 			(NewUnionRule())
// 			break
// 		case '+':
// 			result.AddChild(NewUnionRule())
// 			break
// 		case '*':
// 			rules = append(rules, NewStarGroupRule())
// 		case '(':
// 			starting_parenth = i
// 			break
// 		case ')':

// 			break
// 			// case '{n,m}'
// 			// make a breaker edge
// 			// case '[...]':
// 			// Get All Charcters in list
// 			// Add All Characters individually -- fine since GraphTree
// 			break
// 		// case '[^...]':
// 		case '.':
// 		default:
// 			addRule(NewRuneRule(v))
// 		}
// 	}

// 	return result
// }

// func getGroups(expression string) ([]slice, error) {
// 	slices := make([]slice, 0, 16)
// 	stack := make([]uint, 0, 16)
// 	escape_flag := false

// 	for i, v := range expression {
// 		if v == '\\' {
// 			escape_flag = !escape_flag
// 		} else if escape_flag {
// 			escape_flag = false
// 		} else if v == '(' {
// 			stack = append(stack, uint(i))
// 		} else if v == ')' && len(stack) == 0 {
// 			return []slice{}, errors.New("unmatched ')' in regular expression")
// 		} else if v == ')' {
// 			slices = append(slices, slice{from: stack[len(stack)-1], to: uint(i)})
// 			stack = stack[:len(stack)-1]
// 		}
// 	}

// 	if len(stack) > 0 {
// 		return []slice{}, errors.New("unmatched '(' in regular expression")
// 	}

// 	return slices, nil
// }
