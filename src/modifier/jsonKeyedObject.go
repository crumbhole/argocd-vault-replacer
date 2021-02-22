package modifier

import (
	"encoding/json"
)

type jsonKeyedObjectModifier struct{}

func (_ jsonKeyedObjectModifier) modify(input Kvlist) ([]byte, error) {
	keyArray := make(map[string]string, len(input))
	for _, kv := range input {
		keyArray[string(kv.Key)] = string(kv.Value)
	}
	return json.Marshal(keyArray)
}
