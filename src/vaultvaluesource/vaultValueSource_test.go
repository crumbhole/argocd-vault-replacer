package vaultvaluesource

import (
	"bytes"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"net"
	"testing"
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
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	return ln, client
}

func TestGetValue(t *testing.T) {
	ln, client := createTestVault(t)
	defer ln.Close()

	vs := VaultValueSource{Client: client}

	val, err := vs.GetValue([]byte(`/secret/data/path`), []byte(`foo`))
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if !bytes.Equal(*val, []byte(`hi`)) {
		t.Errorf("/secret/data/path,foo !-> hi, got %s", val)
	}
	val, err = vs.GetValue([]byte(`pa`), []byte(`key`))
	expectedError := `Unexpectedly couldn't find pa~key`
	if err != nil && err.Error() != expectedError {
		t.Errorf("Expecting %s, got %s", expectedError, err)
	}
	if val != nil {
		t.Errorf("pa,key !-> nil, got %s", val)
	}
}
