package substitution

import (
	// "fmt"
	"github.com/joibel/vault-replacer/src/modifier"
	"log"
	"net/url"
	"regexp"
)

type Substitutor struct {
	source ValueSource
}

func unescape(input []byte) []byte {
	result, err := url.QueryUnescape(string(input))
	if err != nil {
		log.Fatal(err)
		return input
	}
	return []byte(result)
}

// Takes the 'dirty' key from the regex and cleans it to the actual key
func getKey(input []byte) []byte {
	reKey := regexp.MustCompile(`^!\s*(.*?)\s*$`)
	keyFound := reKey.FindSubmatch(input)
	if keyFound == nil {
		log.Fatal("Key regex failure")
		return input
	}
	return unescape(keyFound[1])
}

// Takes the 'dirty' modifiers from the regex and turns them into a list
func getModifiers(modifiers []byte) []string {
	reMod := regexp.MustCompile(`^\|\s*(.*?)\s*$`)
	reSplit := regexp.MustCompile(`\s*\|\s*`)
	modsFound := reMod.FindSubmatch(modifiers)
	if modsFound == nil {
		log.Fatal("Mods regex failure")
	}
	return reSplit.Split(string(modsFound[1]), -1)
}

func performModifiers(modifiers []string, input []byte) []byte {
	value := input
	for _, mod := range modifiers {
		value = modifier.Modify(value, mod)
	}
	return value
}

// Swaps a <value:...> for the value from the valuesource
// input should contain no lf/cf
func (s Substitutor) substituteValue(input []byte) []byte {
	reOuter := regexp.MustCompile(`^<\s*vault:\s*([^\!]*[^\s])\s*(\!\s*[^\|]+)?\s*(\|.*)?\s*>$`)
	pathFound := reOuter.FindSubmatch(input)
	// for i, found := range pathFound {
	// 	fmt.Printf("Found[%d] %s\n", i, found)
	// }
	if pathFound != nil {
		if len(pathFound[2]) > 0 {
			path := unescape(pathFound[1])
			key := getKey(pathFound[2])
			modifiers := pathFound[3]
			value := s.source.GetValue(path, key)
			if value == nil {
				log.Printf("Key %v %v lookup failure", path, key)
				return input
			}
			if len(modifiers) > 0 {
				return performModifiers(getModifiers(pathFound[3]), *value)
			}
			return *value
		}
	}
	return input
}
