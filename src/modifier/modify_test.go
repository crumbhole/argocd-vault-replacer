package modifier

import (
	"bytes"
	"testing"
)

func TestNoSuchModifier(t *testing.T) {
	_, err := Modify([]byte(`hi`), "nonsense")
	expectedError := `Invalid modifier nonsense`
	if err.Error() != expectedError {
		t.Errorf("Expecting %s, got %s", expectedError, err)
	}
}

func TestBase64(t *testing.T) {
	tests := map[string]string{
		`Hello`: `SGVsbG8=`,
		`Supersecret
thing`: `U3VwZXJzZWNyZXQKdGhpbmc=`,
	}
	for input, expect := range tests {
		in := []byte(input)
		res, err := Modify(in, "base64")
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}
