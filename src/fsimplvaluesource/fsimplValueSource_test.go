package fsimplvaluesource

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

const (
	urlPrefix = "SECRET_URL_PREFIX"
)

type urlMangleTest struct {
	prefix string
	suffix string
}

func TestURLMangler(t *testing.T) {
	tests := map[string]urlMangleTest{
		`git://github.com/abc/def#branchx`: {
			prefix: `git://github.com/abc/def`,
			suffix: `#branchx`,
		},
		`git+https://github.com/abc/def#branchx`: {
			prefix: `git+https://github.com/abc/def`,
			suffix: `#branchx`,
		},
		`git+https://github.com/abc/def/#branchx`: {
			prefix: `git+https://github.com/abc/def/`,
			suffix: `#branchx`,
		},
		`git+https://github.com/abc/def`: {
			prefix: `git+https://github.com/abc/def`,
			suffix: ``,
		},
		`file://github.com/abc/def#branchx`: {
			prefix: `file://github.com/abc/def#branchx`,
			suffix: ``,
		},
	}

	for url, test := range tests {
		t.Run(url, func(t *testing.T) {
			prefix, suffix := MangleURL(url)
			if prefix != test.prefix {
				t.Errorf(`Unexpected prefix for %s. Got %s, wanted %s`, url, prefix, test.prefix)
			}
			if suffix != test.suffix {
				t.Errorf(`Unexpected suffix for %s. Got %s, wanted %s`, url, suffix, test.suffix)
			}
		})
	}
}

type pathMangleTest struct {
	url  string
	path string
}

func TestPathMangler(t *testing.T) {
	tests := map[pathMangleTest]string{
		{
			url:  `git://github.com/abc/def#branchx`,
			path: `test`,
		}: `test`,
		{
			url:  `https://github.com/abc/def/`,
			path: `/test`,
		}: `/test/`,
		{
			url:  `https://github.com/abc/def/`,
			path: `/test/`,
		}: `/test/`,
	}

	for test, expected := range tests {
		t.Run(fmt.Sprintf("%s:%s", test.url, test.path), func(t *testing.T) {
			result := ManglePathForURL(test.url, test.path)
			if result != expected {
				t.Errorf("Unexpected result for %s %s. Got %s, wanted %s", test.url, test.path, result, expected)
			}
		})
	}
}

func TestGetValueSimple(t *testing.T) {
	vs := New()
	cwd, _ := os.Getwd()
	os.Setenv(urlPrefix, fmt.Sprintf("file://%s/test", cwd))
	val, err := vs.GetValue([]byte(`/test1`), []byte(`foo`))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if !bytes.Equal(*val, []byte(`hi`)) {
		t.Errorf("test1~foo -> hi, got %s", *val)
	}
	val, err = vs.GetValue([]byte(`/pa`), []byte(`key`))
	if err == nil {
		t.Errorf("Expecting error, didn't get one")
	}
	if val != nil {
		t.Errorf("pa,key !-> nil, got %s", val)
	}
}

type FsimplValue struct {
	path string
	key  string
}

func TestParseYaml(t *testing.T) {
	yamlContents := []byte(`
foo:
  bar: orange
  fish:
    cod: trout
`)
	testValues := map[FsimplValue]string{
		{
			path: `foo`,
			key:  `bar`,
		}: `orange`,
		{
			path: `/foo`,
			key:  `bar`,
		}: `orange`,
		{
			path: `foo/`,
			key:  `bar`,
		}: `orange`,
		{
			path: `/foo/`,
			key:  `bar`,
		}: `orange`,
		{
			path: `/foo/fish`,
			key:  `cod`,
		}: `trout`,
	}
	for test, expected := range testValues {
		val, err := parseYaml(yamlContents, []byte(test.path), []byte(test.key))
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		if val == nil {
			t.Fatal("Unexpected nil value")
		}
		if !bytes.Equal(*val, []byte(expected)) {
			t.Errorf("Wanted %s, got %s", expected, *val)
		}
	}
}

func TestGetValues(t *testing.T) {
	vs := New()
	//	cwd, _ := os.Getwd()
	testValues := map[FsimplValue]string{
		{
			path: `/testvalues`,
			key:  `foo`,
		}: `bar`,
		{
			path: `/testvalues/lemon`,
			key:  `fig`,
		}: `banana`,
	}
	sources := map[string]string{
		`git+https://github.com/crumbhole/argocd-vault-replacer.git/#main`:                            `oaQuei1aij`,
		`https://raw.githubusercontent.com/crumbhole/argocd-vault-replacer/main`:                      `oaQuei1aij`, // Same data as git
		`git+https://github.com/crumbhole/argocd-vault-replacer.git/#testdata`:                        `Ooy3phie4o`,
		`git+https://github.com/crumbhole/argocd-vault-replacer.git//testvalues/test.yaml#main`:       `iegeiFe3ae`,
		`https://raw.githubusercontent.com/crumbhole/argocd-vault-replacer/main/testvalues/test.yaml`: `iegeiFe3ae`, // Same data as git
		`git+https://github.com/crumbhole/argocd-vault-replacer.git//testvalues/test.yml#main`:        `vieHuch8yi`,
		`git+https://github.com/crumbhole/argocd-vault-replacer.git//testvalues/test.json#main`:       `bohg2luSai`,
	}

	for url, suffix := range sources {
		for values, prefix := range testValues {
			t.Run(fmt.Sprintf("%s%s¬%s", url, values.path, values.key), func(t *testing.T) {
				os.Setenv(urlPrefix, url)
				val, err := vs.GetValue([]byte(values.path), []byte(values.key))
				if err != nil {
					t.Errorf("Unexpected error %s", err)
				}
				if val == nil {
					t.Fatal("Unexpected nil value")
				}
				expected := fmt.Sprintf("%s-%s", prefix, suffix)
				if !bytes.Equal(*val, []byte(expected)) {
					t.Errorf("%s~%s -> %s, got %s", values.path, values.key, expected, *val)
				}
			})
		}
	}
}
