package substitution

import (
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
	//	fmt.Printf("Unescaped /%s/ to /%s/\n", input, result)
	return []byte(result)
}

// Takes the 'dirty' key from the regex and cleans it to the actual key
func getKey(input []byte) []byte {
	//			fmt.Printf("Key found %d\n", len(pathFound))
	reKey := regexp.MustCompile(`^!\s*(.*?)\s*$`)
	// fmt.Printf("Not found //%v//\n", pathFound[2])
	keyFound := reKey.FindSubmatch(input)
	if keyFound == nil {
		log.Fatal("Key regex failure")
		return input
	}
	return unescape(keyFound[1])
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
			value := s.source.GetValue(path, key)
			if value == nil {
				log.Printf("Key %v %v lookup failure", path, key)
				return input
			}
			return *value
		}
		//		fmt.Printf("No key found %d\n", len(pathFound))
		return pathFound[1]
	}
	return input
}
