package substitution

import (
	"bytes"
	"testing"
)

var mockVs = mockValueSource{values: map[pathKeyTuple][]byte{
	{`/path/to/thing`, `key`}: []byte(`value`),
},
}

func TestMockSource(t *testing.T) {
	val, err := mockVs.GetValue([]byte(`/path/to/thing`), []byte(`key`))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if !bytes.Equal(*val, []byte(`value`)) {
		t.Errorf("/path/to/thing,key !-> value, got %s", val)
	}
}

func TestMockSourceFail(t *testing.T) {
	val, err := mockVs.GetValue([]byte(`pa`), []byte(`key`))
	expectedError := `Couldn't find pa ! key`
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expecting %s, got %s", expectedError, err)
	}
	if val != nil {
		t.Errorf("pa,key !-> nil, got %s", val)
	}
}
