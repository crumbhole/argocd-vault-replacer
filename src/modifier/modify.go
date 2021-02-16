package modifier

import (
	"log"
)

func getModifier(name string) modifier {
	modifiers := map[string]modifier{
		"base64": base64Modifier{},
	}
	println(name)
	if found, ok := modifiers[name]; ok {
		return found
	}
	log.Fatalf("Invalid modifier %s\n", name)
	return nil
}

func Modify(input []byte, name string) []byte {
	modifier := getModifier(name)
	return modifier.modify(input)
}
