package substitution

import (
	"errors"
	"fmt"
	"github.com/crumbhole/argocd-vault-replacer/src/modifier"
	"net/url"
	"regexp"
)

type Substitutor struct {
	source ValueSource
	errs   error
}

func unescape(input []byte) ([]byte, error) {
	result, err := url.QueryUnescape(string(input))
	if err != nil {
		return input, err
	}
	return []byte(result), nil
}

// Takes the 'dirty' key from the regex and cleans it to the actual key
func getKey(input []byte) ([]byte, error) {
	reKey := regexp.MustCompile(`^~\s*(.*?)\s*$`)
	keyFound := reKey.FindSubmatch(input)
	if keyFound == nil {
		return nil, errors.New("Key regex failure")
	}
	return unescape(keyFound[1])
}

// Takes the 'dirty' modifiers from the regex and turns them into a list
func getModifiers(modifiers []byte) ([]string, error) {
	reMod := regexp.MustCompile(`^\|\s*(.*?)\s*$`)
	reSplit := regexp.MustCompile(`\s*\|\s*`)
	modsFound := reMod.FindSubmatch(modifiers)
	if modsFound == nil {
		return nil, errors.New("Mods regex failure")
	}
	return reSplit.Split(string(modsFound[1]), -1), nil
}

func performModifiers(modifiers []string, input []byte) ([]byte, error) {
	value := input
	for _, mod := range modifiers {
		var err error
		value, err = modifier.Modify(value, mod)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

// Swaps a <value:...> for the value from the valuesource
// input should contain no lf/cf
func (s *Substitutor) substituteValueWithError(input []byte) ([]byte, error) {
	reOuter := regexp.MustCompile(`^<\s*vault:\s*([^\~]*[^\s])\s*(\~\s*[^\|]+)?\s*(\|.*)?\s*>$`)
	pathFound := reOuter.FindSubmatch(input)
	if pathFound != nil {
		if len(pathFound[2]) > 0 {
			path, err := unescape(pathFound[1])
			if err != nil {
				return nil, err
			}
			key, err := getKey(pathFound[2])
			if err != nil {
				return nil, err
			}
			modifiers := pathFound[3]
			value, err := s.source.GetValue(path, key)
			if err != nil {
				return nil, err
			}
			if len(modifiers) > 0 {
				modList, err := getModifiers(modifiers)
				if err != nil {
					return nil, err
				}
				return performModifiers(modList, *value)
			}
			return *value, nil
		}
		return nil, errors.New(`Failed to find path for substitution`)
	}
	// We pass through things we can't match at all. They shouldn't arrive here.
	return input, nil
}

// Swaps a <value:...> for the value from the valuesource
// input should contain no lf/cf
func (s *Substitutor) substituteValue(input []byte) []byte {
	res, err := s.substituteValueWithError(input)
	if err != nil {
		if s.errs == nil {
			s.errs = err
		} else {
			s.errs = fmt.Errorf("%s\n%s", s.errs, err)
		}
	}
	return res
}
