package modifier

import (
	"bytes"
	"testing"
)

func TestJsonObjectToList(t *testing.T) {
	// Fragile test - relies on output order of json.Marshal
	tests := map[string]Kvlist{
		`[{"keyname":"key1","valuename":"val1"},{"keyname":"key2","valuename":"val2"}]`: {
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
		},
		`[{"keyname":"key1","valuename":"val1"},{"keyname":"key2","valuename":"val2"},{"keyname":"oink","valuename":"foo"},{"keyname":"sausage","valuename":"bar"}]`: {
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
			{Key: []byte(`oink`), Value: []byte(`foo`)},
			{Key: []byte(`sausage`), Value: []byte(`bar`)},
		},
	}
	modifier := jsonObjectToListModifierGet(`jsonobject2list(keyname,valuename)`)
	for expect, input := range tests {
		res, err := modifier.modifyKvlist(input)
		if err != nil {
			t.Errorf("%v !-> %v, got an error %s", input, expect, err)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%v !-> %v, got %s", input, expect, res)
		}
	}
	modifier = jsonObjectToListModifierGet(`jsonobject2list( keyname   ,    valuename   )`)
	for expect, input := range tests {
		res, err := modifier.modifyKvlist(input)
		if err != nil {
			t.Errorf("%v !-> %v, got an error %s", input, expect, err)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%v !-> %v, got %s", input, expect, res)
		}
	}
}
