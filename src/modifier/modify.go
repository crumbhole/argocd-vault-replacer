package modifier

import (
	"encoding/json"
	"fmt"
	"sort"
)

func textToKvlist(input []byte) (Kvlist, error) {
	list := Kvlist{}
	err := json.Unmarshal(input, &list)
	if err != nil {
		flat := make(map[string]string, 0)
		err = json.Unmarshal(input, &flat)
		if err != nil {
			return nil, err
		}
		for key, val := range flat {
			list = append(list, Kv{Key: []byte(key), Value: []byte(val)})
		}
	}
	// We care that the ordering is stable only
	sort.Slice(list, func(i, j int) bool { return string(list[i].Key) < string(list[j].Key) })
	return list, nil
}

func getModifier(name string) (modifier, error) {
	obj2list := jsonObjectToListModifierGet(name)
	if obj2list != nil {
		return obj2list, nil
	}

	modifiers := map[string]modifier{
		"base64":           base64Modifier{},
		"json2htaccess":    htaccessModifier{},
		"jsonlist":         jsonListModifier{},
		"jsonkeyedobject":  jsonKeyedObjectModifier{},
		"jsonpairedobject": jsonPairedObjectModifier{},
	}
	if found, ok := modifiers[name]; ok {
		return found, nil
	}
	return nil, fmt.Errorf("Invalid modifier %s", name)
}

// Modify takes some input and the name of a modifier to modify that string
// with and returns the changed input.
func Modify(input []byte, name string) ([]byte, error) {
	modifier, err := getModifier(name)
	if err != nil {
		return nil, err
	}
	return modifier.modify(input)
}

func getKVModifier(name string) (kvmodifier, error) {
	obj2list := jsonObjectToListModifierGet(name)
	if obj2list != nil {
		return obj2list, nil
	}

	modifiers := map[string]kvmodifier{
		"valuestext":       valuesTextModifier{},
		"jsonlist":         jsonListModifier{},
		"jsonkeyedobject":  jsonKeyedObjectModifier{},
		"jsonpairedobject": jsonPairedObjectModifier{},
	}
	if found, ok := modifiers[name]; ok {
		return found, nil
	}
	return nil, fmt.Errorf("Invalid modifier %s", name)
}

// ModifyKVList takes some key+values and the name of a modifier to modify that list
// with and returns the changed input.
func ModifyKVList(input Kvlist, name string) ([]byte, error) {
	modifier, err := getKVModifier(name)
	if err != nil {
		return nil, err
	}
	return modifier.modifyKvlist(input)
}
