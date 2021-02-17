package vaultValueSource

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
)

type VaultValueSource struct {
	Client *vault.Client
}

func (m *VaultValueSource) initClient() error {
	if m.Client == nil {
		client, err := vault.NewClient(nil)
		if err != nil {
			return err
		}
		m.Client = client
		err = m.tryKubernetesAuth()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m VaultValueSource) GetValue(path []byte, key []byte) (*[]byte, error) {
	err := m.initClient()
	if err != nil {
		return nil, err
	}
	secret, err := m.Client.Logical().Read(string(path))
	if err != nil {
		return nil, err
	}

	// Joy of casting in go
	switch data := secret.Data["data"].(type) {
	case map[string]interface{}:
		if value, found := data[string(key)]; found {
			switch dataval := value.(type) {
			case string:
				datavalbyte := []byte(dataval)
				return &datavalbyte, nil
			}
		}
	}
	return nil, fmt.Errorf("Couldn't find %s ! %s", path, key)
}
