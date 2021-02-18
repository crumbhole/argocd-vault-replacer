package substitution

import (
	"encoding/base64"
	"regexp"
)

func substituteb64(input []byte) []byte {
	decoded, err := base64.StdEncoding.DecodeString(string(input))
	if err != nil {
		return input
	}
	// This reValue differs from the one in Substitute as it requires a match on
	// the whole thing
	reValue := regexp.MustCompile(`^(<[ \t]*vault:[^\r\n]+?)>$`)
	// We want to stuff |base64 on the end of this
	matchFound := reValue.FindSubmatch(decoded)
	if matchFound != nil {
		return []byte(string(matchFound[1]) + `|base64>`)
	}
	return input
}

// Takes a whole multi-line []byte and finds appropriate subsitutions
func Substitute(input []byte, source ValueSource) ([]byte, error) {
	// First attempt to base64 decode any <vault:> secrets encoded by other
	// tools, such as helm
	reB64Value := regexp.MustCompile(`PHZhdWx0O[^\s]*`)
	debase64input := reB64Value.ReplaceAllFunc(input, substituteb64)

	reValue := regexp.MustCompile(`<[ \t]*vault:[^\r\n]+?>`)
	subst := Substitutor{source: source}
	return reValue.ReplaceAllFunc(debase64input, subst.substituteValue), subst.errs
}
