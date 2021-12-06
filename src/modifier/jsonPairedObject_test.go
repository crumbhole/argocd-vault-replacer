package modifier

import (
	"bytes"
	"testing"
)

func TestJsonPairedObject(t *testing.T) {
	// Fragile test - relies on output order of json.Marshal
	tests := map[string]Kvlist{
		`{"val1":"val2"}`: {
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
		},
		`{"foo":"bar","val1":"val2"}`: {
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
			{Key: []byte(`oink`), Value: []byte(`foo`)},
			{Key: []byte(`sausage`), Value: []byte(`bar`)},
		},
	}
	modifier := jsonPairedObjectModifier{}
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

func TestJsonPairedObjectFail(t *testing.T) {
	tests := []Kvlist{
		{
			{Key: []byte(`key1`), Value: []byte(`val1`)},
			{Key: []byte(`key2`), Value: []byte(`val2`)},
			{Key: []byte(`key3`), Value: []byte(`val3`)},
		},
		{
			{Key: []byte(`key1`), Value: []byte(`val1`)},
		},
	}
	modifier := jsonPairedObjectModifier{}
	for _, input := range tests {
		_, err := modifier.modifyKvlist(input)
		expectedError := `Paired object needs an even number of inputs`
		if err == nil || err.Error() != expectedError {
			t.Errorf("Expecting %s, got %s", expectedError, err)
		}
	}
}
