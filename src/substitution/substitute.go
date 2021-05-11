package substitution

import (
	"encoding/base64"
	"regexp"
)

type Substitutor struct {
	Source ValueSource
	errs   error
}

// Takes a whole multi-line []byte and substitutes with base64 decode
func (s *Substitutor) substitutebase64(input []byte) []byte {
	decoded, err := base64.StdEncoding.DecodeString(string(input))
	// We don't know this is b64, so failure to decode is fine
	if err != nil {
		return input
	}
	// Recurse with the decoded version
	out, err := s.substituteraw(decoded)
	return []byte(base64.StdEncoding.EncodeToString(out))
}

// Takes a whole multi-line []byte and substitutes without base64 decode
func (s *Substitutor) substituteraw(input []byte) ([]byte, error) {
	reValue := regexp.MustCompile(`<[ \t]*vault:[^\r\n]+?>`)
	return reValue.ReplaceAllFunc(input, s.substituteValue), s.errs
}

// Takes a whole multi-line []byte and finds appropriate subsitutions
func (s *Substitutor) Substitute(input []byte) ([]byte, error) {
	// First attempt to base64 decode any <vault:> secrets encoded by other
	// tools, such as helm
	reB64Value := regexp.MustCompile(`[A-Za-z0-9\+\/\=]{10,}`)
	postbase64input := reB64Value.ReplaceAllFunc(input, s.substitutebase64)

	return s.substituteraw(postbase64input)
}
