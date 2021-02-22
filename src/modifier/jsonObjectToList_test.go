package modifier

import (
	"bytes"
	"testing"
)

func TestJsonObjectToList(t *testing.T) {
	// Fragile test - relies on output order of json.Marshal
	tests := map[string]Kvlist{
		`[{"keyname":"key1","valuename":"val1"},{"keyname":"key2","valuename":"val2"}]`: Kvlist{
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
		},
		`[{"keyname":"key1","valuename":"val1"},{"keyname":"key2","valuename":"val2"},{"keyname":"oink","valuename":"foo"},{"keyname":"sausage","valuename":"bar"}]`: Kvlist{
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
			{Key: []byte(`oink`), Value: []byte(`foo`)},
			{Key: []byte(`sausage`), Value: []byte(`bar`)},
		},
	}
	for expect, input := range tests {
		modifier := jsonObjectToListModifierGet(`jsonobject2list(keyname,valuename)`)
		res, err := modifier.modify(input)
		if err != nil {
			t.Errorf("%v !-> %v, got an error %s", input, expect, err)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%v !-> %v, got %s", input, expect, res)
		}
	}
}
