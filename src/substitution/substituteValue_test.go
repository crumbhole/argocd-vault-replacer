package substitution

import (
	"bytes"
	"testing"
)

var substvalVs = mockValueSource{values: map[pathKeyTuple][]byte{
	{`/path/to/thing`, `key`}:           []byte(`value`),
	{`/path/to/thing`, `foo`}:           []byte(`bar`),
	{`/path/to/other`, `nose`}:          []byte(`out`),
	{`/path/to/emoji`, `smile`}:         []byte(`😀`),
	{`/path/to/😀`, `face`}:              []byte(`laugh`),
	{`/path/ /other`, `pear`}:           []byte(`apple`),
	{`/path/ /other`, `ora nge`}:        []byte(`satsu ma`),
	{`/spacepath/ `, `nice`}:            []byte(`time`),
	{`/path/to/thing`, ` leadingspace`}: []byte(`yay`),
	{`/path/>/<`, `<><>`}:               []byte(`pointy`),
},
}
var subst = Substitutor{Source: substvalVs}

func TestBasicFail(t *testing.T) {
	key := []byte(`blah`)
	res, err := subst.substituteValueWithError(key)
	if !bytes.Equal(res, key) {
		t.Errorf("blah !-> blah, got %s", res)
	}
	expectedError := `Failed to find path for substitution`
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expecting %s, got %s", expectedError, err)
	}
}

// func TestBasicSuccess(t *testing.T) {
// 	key := []byte(`<secret:/path/to/thing>`)
// 	res := subst.substituteValue(key)
// 	if !bytes.Equal(res, []byte(`/path/to/thing`)) {
// 		t.Errorf("<vault:/path/to/thing> !-> /path/to/thing, got %s", res)
// 	}
// }

func TestManyGood(t *testing.T) {
	tests := map[string]string{
		`<secret:/path/to/thing~key>`:                   `value`,
		`<secret:/path/to/thing/~key>`:                  `value`,
		`<secret:/path/to/thing~foo>`:                   `bar`,
		`< secret:/path/to/thing~key>`:                  `value`,
		`<secret: /path/to/thing~key>`:                  `value`,
		`<secret:/path/to/thing ~key>`:                  `value`,
		`<secret:/path/to/thing~ key>`:                  `value`,
		`<secret:/path/to/thing~key >`:                  `value`,
		`< secret: /path/to/thing ~ key >`:              `value`,
		`<  secret:  /path/to/thing  ~  key  >`:         `value`,
		`<secret:/path/to/other~nose>`:                  `out`,
		`<secret:/path/to/😀 ~ face >`:                   `laugh`,
		`<secret:/path/to/emoji ~ smile >`:              `😀`,
		`<secret:/path/ /other~pear >`:                  `apple`,
		`<secret:/path/%20/other~pear >`:                `apple`,
		`<secret:/path/ /other~ora nge >`:               `satsu ma`,
		`<secret:/path/%20/other~ora%20nge >`:           `satsu ma`,
		`<secret:/spacepath/%20 ~ nice >`:               `time`,
		`<secret:/path/to/thing ~ %20leadingspace >`:    `yay`,
		`<secret:/path/%3E/%3c/~%3c%3e%3c%3e>`:          `pointy`,
		`<secret:/path/to/thing~key|base64>`:            `dmFsdWU=`,
		`<secret:/path/to/thing~key | base64  >`:        `dmFsdWU=`,
		`<secret:/path/to/thing~key | base64  |base64>`: `ZG1Gc2RXVT0=`,

		`<vault:/path/to/thing~key>`:                   `value`,
		`<vault:/path/to/thing/~key>`:                  `value`,
		`<vault:/path/to/thing~foo>`:                   `bar`,
		`< vault:/path/to/thing~key>`:                  `value`,
		`<vault: /path/to/thing~key>`:                  `value`,
		`<vault:/path/to/thing ~key>`:                  `value`,
		`<vault:/path/to/thing~ key>`:                  `value`,
		`<vault:/path/to/thing~key >`:                  `value`,
		`< vault: /path/to/thing ~ key >`:              `value`,
		`<  vault:  /path/to/thing  ~  key  >`:         `value`,
		`<vault:/path/to/other~nose>`:                  `out`,
		`<vault:/path/to/😀 ~ face >`:                   `laugh`,
		`<vault:/path/to/emoji ~ smile >`:              `😀`,
		`<vault:/path/ /other~pear >`:                  `apple`,
		`<vault:/path/%20/other~pear >`:                `apple`,
		`<vault:/path/ /other~ora nge >`:               `satsu ma`,
		`<vault:/path/%20/other~ora%20nge >`:           `satsu ma`,
		`<vault:/spacepath/%20 ~ nice >`:               `time`,
		`<vault:/path/to/thing ~ %20leadingspace >`:    `yay`,
		`<vault:/path/%3E/%3c/~%3c%3e%3c%3e>`:          `pointy`,
		`<vault:/path/to/thing~key|base64>`:            `dmFsdWU=`,
		`<vault:/path/to/thing~key | base64  >`:        `dmFsdWU=`,
		`<vault:/path/to/thing~key | base64  |base64>`: `ZG1Gc2RXVT0=`,
	}
	for input, expect := range tests {
		in := []byte(input)
		res, err := subst.substituteValueWithError(in)
		if err != nil {
			t.Errorf("%s !-> %v, got an error %s", in, expect, err)
		}
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}

func TestManyBad(t *testing.T) {
	tests := []string{
		`<secret:/path/to/thing~key`,
		`secret:/path/to/thing~key>`,
		`<ecret:/path/to/thing~key>`,
		`<secret/path/to/thing~key>`,

		`<vault:/path/to/thing~key`,
		`vault:/path/to/thing~key>`,
		`<ault:/path/to/thing~key>`,
		`<vault/path/to/thing~key>`,
	}
	for _, input := range tests {
		in := []byte(input)
		res, err := subst.substituteValueWithError(in)
		if err != nil {
			t.Errorf("want %s untouched, got an error %s", in, err)
		}
		if !bytes.Equal(res, in) {
			t.Errorf("want %s untouched but got %s", input, res)
		}
	}
}

func TestBadSubst(t *testing.T) {
	tests := []string{
		`<vault:/path/to/thing~key|nonsense>`,
	}
	for _, input := range tests {
		in := []byte(input)
		_, err := subst.substituteValueWithError(in)
		expectedError := `Invalid modifier nonsense`
		if err != nil && err.Error() != expectedError {
			t.Errorf("Expecting %s, got %s", expectedError, err)
		}
	}
}
