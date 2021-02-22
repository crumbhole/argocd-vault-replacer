package modifier

import (
	"fmt"
)

func getModifier(name string) (modifier, error) {
	modifiers := map[string]modifier{
		"base64": base64Modifier{},
	}
	if found, ok := modifiers[name]; ok {
		return found, nil
	}
	return nil, fmt.Errorf("Invalid modifier %s", name)
}

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

func ModifyKVList(input Kvlist, name string) ([]byte, error) {
	modifier, err := getKVModifier(name)
	if err != nil {
		return nil, err
	}
	return modifier.modify(input)
}
