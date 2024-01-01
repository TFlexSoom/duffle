package lexer

import (
	decision "github.com/tflexsoom/duffle/internal/dfa"
)

type EitherStringByte struct {
	StringVal string
	ByteVal   []byte
}

type Token struct {
	TokenId int
	Val     EitherStringByte
}

type Lexer struct {
	tokenIdToName map[int]string
	decisionTree  *decision.DfaGraph[int]
}

type LexerRule struct {
	Name   string
	Regexp string
}

func FromRules(rules []LexerRule) Lexer {
	rulesLen := len(rules)
	tokenIdToName := make(map[int]string, rulesLen)
	decisionTree := decision.NewDfaGraph[int]()

	for i, val := range rules {
		tokenIdToName[i] = val.Name
		decision.AppendRegexp(val.Regexp, i, decisionTree)
	}

	return Lexer{
		tokenIdToName: tokenIdToName,
		decisionTree:  decisionTree,
	}
}
