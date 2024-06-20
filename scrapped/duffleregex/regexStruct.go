package duffleregex

type ruleFlag uint8

const (
	RULE_FLAG_NONE       uint8 = 0
	RULE_FLAG_TO_SUCCESS uint8 = 1
	RULE_FLAG_TO_FAIL    uint8 = 2
)

type rule struct {
	contains []rune
	flag     ruleFlag
}

type state struct {
	regVal      rune
	transitions []rule
}

type DuffleRegex struct {
	input  []rune
	regex  []state
	search []search
}

/////////////////////////
