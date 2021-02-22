package modifier

import (
	"bytes"
	"testing"
)

func TestJsonList(t *testing.T) {
	tests := map[string]Kvlist{
		`["val1","val2"]`: Kvlist{
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
		},
	}
	for expect, input := range tests {
		modifier := jsonListModifier{}
		res, err := modifier.modify(input)
		if err != nil {
			t.Errorf("%v !-> %v, got an error %s", input, expect, err)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%v !-> %v, got %s", input, expect, res)
		}
	}
}
