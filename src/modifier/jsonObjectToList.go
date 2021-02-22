package modifier

import (
	"encoding/json"
	"regexp"
)

// This is special as it takes parameters

func jsonObjectToListModifierGet(call string) kvmodifier {
	// TODO: Some helpful to the user error handling here when parsing fails
	reParams := regexp.MustCompile(`jsonobject2list\(([^\,\)\s)]+?)\,([^\,\)\s)]+?)\)`)
	paramsFound := reParams.FindStringSubmatch(call)
	if paramsFound == nil {
		return nil
	}
	modifier := jsonObjectToListModifier{
		keyname:   string(paramsFound[1]),
		valuename: string(paramsFound[2]),
	}
	return modifier
}

type jsonObjectToListModifier struct {
	keyname   string
	valuename string
}

func (mod jsonObjectToListModifier) modify(input Kvlist) ([]byte, error) {
	list := make([]map[string]string, 0)
	for _, kv := range input {
		newObj := make(map[string]string)
		newObj[mod.keyname] = string(kv.Key)
		newObj[mod.valuename] = string(kv.Value)
		list = append(list, newObj)
	}
	return json.Marshal(list)
}
