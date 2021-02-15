package substitution

import (
	"regexp"
)

// Takes a whole multi-line []byte and finds appropriate subsitutions
func Substitute(input []byte, source ValueSource) []byte {
	reValue := regexp.MustCompile(`<[ \t]*vault:[^\r\n]+?>`)
	subst := Substitutor{source: source}
	return reValue.ReplaceAllFunc(input, subst.substituteValue)
}
