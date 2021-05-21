package substitution

import (
	"bytes"
	"testing"
)

var subst_vs = mockValueSource{values: map[pathKeyTuple][]byte{
	{`/path/to/thing`, `foo`}:    []byte(`bar`),
	{`/path/to/thing`, `frog`}:   []byte(`wallop`),
	{`/path/to/thing`, `really`}: []byte(`nice`),
	{`/spacepath/ `, `nice`}:     []byte(`time`),
},
}

func TestStringSubst(t *testing.T) {
	tests := map[string]string{
		`Hello, we're looking for foo <secret:/path/to/thing~foo> to be here`: `Hello, we're looking for foo bar to be here`,

		`Hello, we're looking for foo<secret:/path/to/thing~foo> to be here`: `Hello, we're looking for foobar to be here`,

		`Hi foo<secret:/path/to/thing~foo> <secret:/spacepath/%20~nice>.`: `Hi foobar time.`,

		`Hi foo<secret:/path/to/thing
~foo> <secret:/spacepath/%20~nice>.`: `Hi foo<secret:/path/to/thing
~foo> time.`,
		`Hello, my secret is <secret:/path/to/thing~foo|base64>.`:              `Hello, my secret is YmFy.`,
		`Hi foo<secret:/path/to/thing~foo~frog> <secret:/spacepath/%20~nice>.`: `Hi foobarwallop time.`,

		`Hello, we're looking for foo <vault:/path/to/thing~foo> to be here`: `Hello, we're looking for foo bar to be here`,

		`Hello, we're looking for foo<vault:/path/to/thing~foo> to be here`: `Hello, we're looking for foobar to be here`,

		`Hi foo<vault:/path/to/thing~foo> <vault:/spacepath/%20~nice>.`: `Hi foobar time.`,

		`Hi foo<vault:/path/to/thing
~foo> <vault:/spacepath/%20~nice>.`: `Hi foo<vault:/path/to/thing
~foo> time.`,
		`Hello, my secret is <vault:/path/to/thing~foo|base64>.`:             `Hello, my secret is YmFy.`,
		`Hi foo<vault:/path/to/thing~foo~frog> <vault:/spacepath/%20~nice>.`: `Hi foobarwallop time.`,
	}
	for input, expect := range tests {
		in := []byte(input)
		subst := Substitutor{Source: subst_vs}
		res, errs := subst.Substitute(in)
		if errs != nil {
			t.Errorf("Got unexpected errors in substitute test %s", errs)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}

func TestStringSubstB64(t *testing.T) {
	tests := map[string]string{
		`foo PHNlY3JldDovcGF0aC90by90aGluZ35mb28+ to be here`:                                                                          `foo YmFy to be here`,
		`foo "PHNlY3JldDovcGF0aC90by90aGluZ35mb28+" to be here`:                                                                        `foo "YmFy" to be here`,
		`foo 'PHNlY3JldDovcGF0aC90by90aGluZ35mb28+' to be here`:                                                                        `foo 'YmFy' to be here`,
		`fooPHNlY3JldDovcGF0aC90by90aGluZ35mb28+ to be here`:                                                                           `fooPHNlY3JldDovcGF0aC90by90aGluZ35mb28+ to be here`,
		`fooPHNlY3JldDovcGF0aC90by90aGluZ35mb28+to be here`:                                                                            `fooPHNlY3JldDovcGF0aC90by90aGluZ35mb28+to be here`,
		`VGhpcyBpcyBhIG1peGVkIHVwIDxzZWNyZXQ6L3BhdGgvdG8vdGhpbmd+Zm9vPiB0aGluZyA8c2VjcmV0Oi9zcGFjZXBhdGgvJTIwfm5pY2U+IGluIGJhc2U2NA==`: `VGhpcyBpcyBhIG1peGVkIHVwIGJhciB0aGluZyB0aW1lIGluIGJhc2U2NA==`,

		`foo PHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4= to be here`:                                                                      `foo YmFy to be here`,
		`foo "PHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=" to be here`:                                                                    `foo "YmFy" to be here`,
		`foo 'PHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=' to be here`:                                                                    `foo 'YmFy' to be here`,
		`fooPHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4= to be here`:                                                                       `fooPHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4= to be here`,
		`fooPHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=to be here`:                                                                        `fooPHZhdWx0Oi9wYXRoL3RvL3RoaW5nfmZvbz4=to be here`,
		`VGhpcyBpcyBhIG1peGVkIHVwIDx2YXVsdDovcGF0aC90by90aGluZ35mb28+IHRoaW5nIDx2YXVsdDovc3BhY2VwYXRoLyUyMH5uaWNlPiBpbiBiYXNlNjQ=`: `VGhpcyBpcyBhIG1peGVkIHVwIGJhciB0aGluZyB0aW1lIGluIGJhc2U2NA==`,
	}
	for input, expect := range tests {
		in := []byte(input)
		subst := Substitutor{Source: subst_vs}
		res, errs := subst.Substitute(in)
		if errs != nil {
			t.Errorf("Got unexpected errors in substitute test %s", errs)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}

func TestModifiers(t *testing.T) {
	// Fragile ordering
	tests := map[string]string{
		`JSON <secret:/path/to/thing~foo~frog|jsonlist>.`:                                                `JSON ["bar","wallop"].`,
		`JSON <secret:/path/to/thing~frog~foo|jsonlist>.`:                                                `JSON ["wallop","bar"].`,
		`JSON <secret:/path/to/thing~frog~foo|jsonkeyedobject>.`:                                         `JSON {"foo":"bar","frog":"wallop"}.`,
		`JSON <secret:/path/to/thing~frog~foo|jsonpairedobject>.`:                                        `JSON {"wallop":"bar"}.`,
		`JSON <secret:/path/to/thing~foo~frog|jsonpairedobject>.`:                                        `JSON {"bar":"wallop"}.`,
		`JSON <secret:/path/to/thing~foo~frog|jsonpairedobject|jsonobject2list(key1,key2)>.`:             `JSON [{"key1":"bar","key2":"wallop"}].`,
		`JSON <secret:/path/to/thing~foo~frog~frog~really|jsonpairedobject|jsonobject2list(key1,key2)>.`: `JSON [{"key1":"bar","key2":"wallop"},{"key1":"wallop","key2":"nice"}].`,
		`JSON <secret:/path/to/thing~foo~frog|jsonkeyedobject|jsonobject2list( key1 , key2 )>.`:          `JSON [{"key1":"foo","key2":"bar"},{"key1":"frog","key2":"wallop"}].`,

		`JSON <vault:/path/to/thing~foo~frog|jsonlist>.`:                                                `JSON ["bar","wallop"].`,
		`JSON <vault:/path/to/thing~frog~foo|jsonlist>.`:                                                `JSON ["wallop","bar"].`,
		`JSON <vault:/path/to/thing~frog~foo|jsonkeyedobject>.`:                                         `JSON {"foo":"bar","frog":"wallop"}.`,
		`JSON <vault:/path/to/thing~frog~foo|jsonpairedobject>.`:                                        `JSON {"wallop":"bar"}.`,
		`JSON <vault:/path/to/thing~foo~frog|jsonpairedobject>.`:                                        `JSON {"bar":"wallop"}.`,
		`JSON <vault:/path/to/thing~foo~frog|jsonpairedobject|jsonobject2list(key1,key2)>.`:             `JSON [{"key1":"bar","key2":"wallop"}].`,
		`JSON <vault:/path/to/thing~foo~frog~frog~really|jsonpairedobject|jsonobject2list(key1,key2)>.`: `JSON [{"key1":"bar","key2":"wallop"},{"key1":"wallop","key2":"nice"}].`,
		`JSON <vault:/path/to/thing~foo~frog|jsonkeyedobject|jsonobject2list( key1 , key2 )>.`:          `JSON [{"key1":"foo","key2":"bar"},{"key1":"frog","key2":"wallop"}].`,
	}
	for input, expect := range tests {
		in := []byte(input)
		subst := Substitutor{Source: subst_vs}
		res, errs := subst.Substitute(in)
		if errs != nil {
			t.Errorf("Got unexpected errors in substitute test %s", errs)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}
