package modifier

import (
	"bytes"
	"testing"
)

func TestBase64(t *testing.T) {
	tests := map[string]string{
		`Hello`: `SGVsbG8=`,
		`Supersecret
thing`: `U3VwZXJzZWNyZXQKdGhpbmc=`,
	}
	for input, expect := range tests {
		in := []byte(input)
		res := Modify(in, "base64")
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}

}
