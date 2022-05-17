package substitution

import (
	"bytes"
	"encoding/base64"
	"regexp"
)

// Substitutor is acting like a class to hold the information to perform substitution on some data
// and collect the errors during that substitution.
type Substitutor struct {
	Source ValueSource
	errs   error
}

// substitutebase64 takes a whole multi-line []byte and substitutes with base64 decode
func (s *Substitutor) substitutebase64(input []byte) []byte {
	decoded, err := base64.StdEncoding.DecodeString(string(input))
	// We don't know this is b64, so failure to decode is fine
	if err != nil {
		return input
	}
	// Recurse with the decoded version
	out, _ := s.substituteraw(decoded)

	// Fast exit if subsititution has done nothing
	// Works around a bug where we have something does not b64 encode->decode
	// without without changing, for things that don't need substution
	if bytes.Equal(out, decoded) {
		return input
	}
	return []byte(base64.StdEncoding.EncodeToString(out))
}

// substituteraw takes a whole multi-line []byte and substitutes without base64 decode
func (s *Substitutor) substituteraw(input []byte) ([]byte, error) {
	reValue := regexp.MustCompile(`<[ \t]*(secret|vault):[^\r\n]+?>`)
	return reValue.ReplaceAllFunc(input, s.substituteValue), s.errs
}

// Substitute takes a whole multi-line []byte and finds appropriate subsitutions
func (s *Substitutor) Substitute(input []byte) ([]byte, error) {
	// First attempt to base64 decode any <vault:> secrets encoded by other
	// tools, such as helm
	reB64Value := regexp.MustCompile(`[A-Za-z0-9\+\/\=]{10,}`)
	postbase64input := reB64Value.ReplaceAllFunc(input, s.substitutebase64)

	return s.substituteraw(postbase64input)
}
