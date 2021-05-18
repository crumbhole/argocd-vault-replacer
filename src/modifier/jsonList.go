package modifier

import (
	"encoding/json"
)

type jsonListModifier struct {
}

func (this jsonListModifier) modify(input []byte) ([]byte, error) {
	list, err := textToKvlist(input)
	if err != nil {
		return nil, err
	}
	return this.modifyKvlist(list)
}

func (_ jsonListModifier) modifyKvlist(input Kvlist) ([]byte, error) {
	keyArray := make([]string, 0)
	for _, value := range input {
		keyArray = append(keyArray, string(value.Value))
	}
	return json.Marshal(keyArray)
}
