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

type mangleTest struct {
	prefix string
	suffix string
}

func TestUrlMangler(t *testing.T) {
	tests := map[string]mangleTest{
		"git://github.com/abc/def#branchx": {
			prefix: `git://github.com/abc/def`,
			suffix: `#branchx`,
		},
		"git+https://github.com/abc/def#branchx": {
			prefix: `git+https://github.com/abc/def`,
			suffix: `#branchx`,
		},
		"git+https://github.com/abc/def": {
			prefix: `git+https://github.com/abc/def`,
			suffix: ``,
		},
		"file://github.com/abc/def#branchx": {
			prefix: `file://github.com/abc/def#branchx`,
			suffix: ``,
		},
	}

	for url, test := range tests {
		t.Run(url, func(t *testing.T) {
			prefix, suffix := MangleUrl(url)
			if prefix != test.prefix {
				t.Errorf("Unexpected prefix for %s. Got %s, wanted %s", url, prefix, test.prefix)
			}
			if suffix != test.suffix {
				t.Errorf("Unexpected suffix for %s. Got %s, wanted %s", url, suffix, test.suffix)
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

type FsimplTest struct {
	Url    string
	Path   string
	Key    string
	Result string
}

func TestGetValues(t *testing.T) {
	vs := New()
	//	cwd, _ := os.Getwd()
	tests := map[string]FsimplTest{
		"Github over https": {
			Url:    "git+https://github.com/crumbhole/argocd-vault-replacer.git/#main",
			Path:   "/testvalues",
			Key:    "foo",
			Result: "bar",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			os.Setenv(urlPrefix, test.Url)
			val, err := vs.GetValue([]byte(test.Path), []byte(test.Key))
			if err != nil {
				t.Errorf("Unexpected error %s", err)
			}
			if !bytes.Equal(*val, []byte(test.Result)) {
				t.Errorf("%s~%s -> %s, got %s", test.Path, test.Key, test.Result, *val)
			}
		})
	}
}
