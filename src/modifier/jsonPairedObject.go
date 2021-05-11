package modifier

import (
	"encoding/json"
	"errors"
)

type jsonPairedObjectModifier struct{}

func (_ jsonPairedObjectModifier) modify(input Kvlist) ([]byte, error) {
	if len(input)%2 != 0 {
		return nil, errors.New(`Paired object needs an even number of inputs`)
	}
	keyArray := make(map[string]string, len(input))
	for index := 0; index < len(input); index += 2 {
		keyArray[string(input[index].Value)] = string(input[index+1].Value)
	}
	return json.Marshal(keyArray)
}
