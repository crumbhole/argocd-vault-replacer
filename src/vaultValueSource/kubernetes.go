package vaultValueSource

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	serviceAccountFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	roleEnv            = "VAULT_ROLE"
	defaultRole        = "argocd"
	authPathEnv        = "VAULT_AUTH_PATH"
	defaultAuthPath    = "/auth/kubernetes/login/"
)

// readJWT reads the JWT data for the Agent to submit to Vault. The default is
// to read the JWT from the default service account location, defined by the
// constant serviceAccountFile. In normal use k.jwtData is nil at invocation and
// the method falls back to reading the token path with os.Open, opening a file
// from either the default location or from the token_path path specified in
// configuration.
func readJWT() (string, error) {
	// load configured token path if set, default to serviceAccountFile
	tokenFilePath := serviceAccountFile

	f, err := os.Open(tokenFilePath)
	if err != nil {
		log.Printf("Kubernetes authentication - no secret found %v", err)
		return "", nil
	}
	defer f.Close()

	contentBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contentBytes)), nil
}

func getVaultRole() string {
	if val, ok := os.LookupEnv(roleEnv); ok {
		return val
	}
	return defaultRole
}

func getVaultAuthPath() string {
	if val, ok := os.LookupEnv(authPathEnv); ok {
		return fmt.Sprintf("/auth/%s/login/", val)
	}
	return defaultAuthPath
}

func (m *VaultValueSource) tryKubernetesAuth() error {
	jwt, err := readJWT()
	if err != nil {
		return err
	}
	if jwt == "" {
		return nil
	}
	secret, err := m.Client.Logical().Write("/auth/kubernetes/login/", map[string]interface{}{
		"role": getVaultRole(),
		"jwt":  jwt,
	})
	if err != nil {
		return err
	}
	m.Client.SetToken(secret.Auth.ClientToken)
	return nil
}
