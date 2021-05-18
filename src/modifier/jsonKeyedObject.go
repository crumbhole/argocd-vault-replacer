package modifier

import (
	"encoding/json"
)

type jsonKeyedObjectModifier struct {
}

func (this jsonKeyedObjectModifier) modify(input []byte) ([]byte, error) {
	list, err := textToKvlist(input)
	if err != nil {
		return nil, err
	}
	return this.modifyKvlist(list)
}

func (_ jsonKeyedObjectModifier) modifyKvlist(input Kvlist) ([]byte, error) {
	keyArray := make(map[string]string, len(input))
	for _, kv := range input {
		keyArray[string(kv.Key)] = string(kv.Value)
	}
	return json.Marshal(keyArray)
}
