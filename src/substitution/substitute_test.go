package substitution

import (
	"bytes"
	"testing"
)

var subst_vs = mockValueSource{values: map[pathKeyTuple][]byte{
	{`/path/to/thing`, `foo`}: []byte(`bar`),
	{`/spacepath/ `, `nice`}:  []byte(`time`),
},
}

func TestStringSubst(t *testing.T) {
	tests := map[string]string{
		`Hello, we're looking for foo <vault:/path/to/thing!foo> to be here`: `Hello, we're looking for foo bar to be here`,

		`Hello, we're looking for foo<vault:/path/to/thing!foo> to be here`: `Hello, we're looking for foobar to be here`,

		`Hi foo<vault:/path/to/thing!foo> <vault:/spacepath/%20!nice>.`: `Hi foobar time.`,

		`Hi foo<vault:/path/to/thing
!foo> <vault:/spacepath/%20!nice>.`: `Hi foo<vault:/path/to/thing
!foo> time.`,
		`Hello, my secret is <vault:/path/to/thing!foo|base64>.`: `Hello, my secret is YmFy.`,
	}
	for input, expect := range tests {
		in := []byte(input)
		res := Substitute(in, subst_vs)
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}
