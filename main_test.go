package main

import (
	"bytes"
	"fmt"
	"github.com/crumbhole/argocd-vault-replacer/src/vaultvaluesource"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"io"
	"io/ioutil"
	"net"
	"os"
	"testing"
)

const (
	testsPath = "test/"
)

func createTestVault(t *testing.T) (net.Listener, *api.Client) {
	t.Helper()

	// Create an in-memory, unsealed core (the "backend", if you will).
	core, keyShares, rootToken := vault.TestCoreUnsealed(t)
	_ = keyShares

	// Start an HTTP server for the core.
	ln, addr := http.TestServer(t, core)

	// Create a client that talks to the server, initially authenticating with
	// the root token.
	conf := api.DefaultConfig()
	conf.Address = addr

	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(rootToken)

	// Setup required secrets, policies, etc.
	_, err = client.Logical().Write("secret/data/path", map[string]interface{}{
		"data": map[string]interface{}{
			"foo": "hi",
			"bar": "example",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	return ln, client
}

func checkDir(t *testing.T, s scanner, path string) error {
	oldstdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	outC := make(chan string)

	err := s.scanDir(path)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	w.Close()
	os.Stdout = oldstdout
	if err != nil {
		return err
	}
	out := <-outC

	expected, err := ioutil.ReadFile(path + "/expected.txt")
	if err != nil {
		return err
	}
	if out != string(expected) {
		return fmt.Errorf("Expected %s and got %s", expected, out)
	}
	return nil
}

// Finds directories under ./test and substitutes all the .yaml/.ymls
// against the above vault, expecting to see expected.txt as the output
func TestDirectories(t *testing.T) {
	ln, client := createTestVault(t)
	defer ln.Close()

	dirs, err := ioutil.ReadDir(testsPath)
	if err != nil {
		t.Error(err)
	}
	s := scanner{source: vaultvaluesource.VaultValueSource{Client: client}}

	for _, d := range dirs {
		t.Run(d.Name(), func(t *testing.T) {
			t.Logf("Testing dir %s", testsPath+d.Name())
			err := checkDir(t, s, testsPath+d.Name())
			if err != nil {
				t.Error(err)
			}
		})
	}
}
