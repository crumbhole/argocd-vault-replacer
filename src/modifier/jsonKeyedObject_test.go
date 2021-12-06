package modifier

import (
	"bytes"
	"testing"
)

func TestJsonKeyedObject(t *testing.T) {
	// Fragile test - relies on output order of json.Marshal
	tests := map[string]Kvlist{
		`{"key1":"val1","key2":"val2"}`: {
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
		},
		`{"key1":"val1","key2":"val2","oink":"foo","sausage":"bar"}`: {
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
			{Key: []byte(`oink`), Value: []byte(`foo`)},
			{Key: []byte(`sausage`), Value: []byte(`bar`)},
		},
	}
	modifier := jsonKeyedObjectModifier{}
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
