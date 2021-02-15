package substitution

import (
	"bytes"
	"testing"
)

var vs = mockValueSource{values: map[pathKeyTuple][]byte{
	{`/path/to/thing`, `key`}: []byte(`value`),
	{`/path/to/thing`, `foo`}: []byte(`bar`),
	{`/path/to/other`, `nose`}: []byte(`out`),
},
}

func TestBasicFail(t *testing.T) {
	key := []byte(`blah`)
	res := substituteValue(key, vs)
	if !bytes.Equal(res, key) {
		t.Errorf("blah !-> blah, got %s", res)
	}
}

func TestBasicSuccess(t *testing.T) {
	key := []byte(`<value:/path/to/thing>`)
	res := substituteValue(key, vs)
	if !bytes.Equal(res, []byte(`/path/to/thing`)) {
		t.Errorf("<value:/path/to/thing> !-> /path/to/thing, got %s", res)
	}
}

func TestManyGood(t *testing.T) {
	tests := map[string]string{
		`<value:/path/to/thing!key>`:`value`,
		`<value:/path/to/thing/!key>`:`value`,
		`<value:/path/to/thing!foo>`:`bar`,
		`< value:/path/to/thing!key>`:`value`,
		`<value: /path/to/thing!key>`:`value`,
		`<value:/path/to/thing !key>`:`value`,
		`<value:/path/to/thing! key>`:`value`,
		`<value:/path/to/thing!key >`:`value`,
		`< value: /path/to/thing ! key >`:`value`,
		`<  value:  /path/to/thing  !  key  >`:`value`,
		`<value:/path/to/other!nose>`:`out`,
	}
	for input, expect := range tests{
		in := []byte(input)
		res := substituteValue(in, vs)
		if !bytes.Equal(res, []byte(expect)) {
			t.Errorf("%s !-> %v, got %s", in, expect, res)
		}
	}
}

func TestManyBad(t *testing.T) {
	tests := []string{
		`<value:/path/to/thing!key`,
		`value:/path/to/thing!key>`,
		`<alue:/path/to/thing!key>`,
		`<value/path/to/thing!key>`,
	}
	for _, input := range tests{
		in := []byte(input)
		res := substituteValue(in, vs)
		if !bytes.Equal(res, in) {
			t.Errorf("want %s untouched but got %s", input, res)
		}
	}
}
