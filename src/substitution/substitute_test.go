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
		`Hello, we're looking for foo <vault:/path/to/thing~foo> to be here`: `Hello, we're looking for foo bar to be here`,

		`Hello, we're looking for foo<vault:/path/to/thing~foo> to be here`: `Hello, we're looking for foobar to be here`,

		`Hi foo<vault:/path/to/thing~foo> <vault:/spacepath/%20~nice>.`: `Hi foobar time.`,

		`Hi foo<vault:/path/to/thing
~foo> <vault:/spacepath/%20~nice>.`: `Hi foo<vault:/path/to/thing
~foo> time.`,
		`Hello, my secret is <vault:/path/to/thing~foo|base64>.`: `Hello, my secret is YmFy.`,
	}
	for input, expect := range tests {
		in := []byte(input)
		res, errs := Substitute(in, subst_vs)
		if errs != nil {
			t.Errorf("Got unexpected errors in substitute test %s", errs)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}

func TestB64Decode(t *testing.T) {
	tests := map[string]string{
		`Hello`:                                    `Hello`,
		`<vault:/path/to/thing~foo>`:               `<vault:/path/to/thing~foo>`,
		`PHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=`:     `<vault:/path/to/thing~foo|base64>`,
		`PHZhdWx0Oi9wYXRoL3RvL3RoaW5nIH4gZm9vID4=`: `<vault:/path/to/thing ~ foo |base64>`,
		`a PHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4= b`: `a PHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4= b`,
		`aPHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=b`:   `aPHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=b`,
	}
	for input, expect := range tests {
		in := []byte(input)
		res := substituteb64(in)
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}

func TestStringSubstB64(t *testing.T) {
	tests := map[string]string{
		`foo PHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4= to be here`:   `foo YmFy to be here`,
		`foo "PHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=" to be here`: `foo "YmFy" to be here`,
		`foo 'PHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=' to be here`: `foo 'YmFy' to be here`,
		`fooPHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4= to be here`:    `fooYmFy to be here`,
		`fooPHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=to be here`:     `fooPHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=to be here`,
	}
	for input, expect := range tests {
		in := []byte(input)
		res, errs := Substitute(in, subst_vs)
		if errs != nil {
			t.Errorf("Got unexpected errors in substitute test %s", errs)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}
