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
