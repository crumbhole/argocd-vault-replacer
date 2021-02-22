package modifier

import (
	"encoding/json"
)

type jsonListModifier struct{}

func (_ jsonListModifier) modify(input Kvlist) ([]byte, error) {
	keyArray := make([]string, 0)
	for _, value := range input {
		keyArray = append(keyArray, string(value.Value))
	}
	return json.Marshal(keyArray)
}
