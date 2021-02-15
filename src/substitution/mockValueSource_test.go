package substitution

import (
	"bytes"
	"testing"
)

var mock_vs = mockValueSource{values: map[pathKeyTuple][]byte{
	{`/path/to/thing`, `key`}: []byte(`value`),
},
}

func TestMockSource(t *testing.T) {
	val := mock_vs.GetValue([]byte(`/path/to/thing`), []byte(`key`))
	if !bytes.Equal(*val, []byte(`value`)) {
		t.Errorf("/path/to/thing,key !-> value, got %s", val)
	}
}

func TestMockSourceFail(t *testing.T) {
	val := mock_vs.GetValue([]byte(`pa`), []byte(`key`))
	if val != nil {
		t.Errorf("pa,key !-> nil, got %s", val)
	}
}
