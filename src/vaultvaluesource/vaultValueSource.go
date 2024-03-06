package vaultvaluesource

import (
	"fmt"
	"github.com/openbao/openbao/api"
)

// VaultValueSource is a value source getting values from hashicorp vault
type VaultValueSource struct {
	Client *api.Client
}

func (m *VaultValueSource) initClient() error {
	if m.Client == nil {
		client, err := api.NewClient(nil)
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

// GetValue returns a value from a path+key in hashicorp vault or null if it doesn't exist
func (m VaultValueSource) GetValue(path []byte, key []byte) (*[]byte, error) {
	err := m.initClient()
	if err != nil {
		return nil, err
	}
	secret, err := m.Client.Logical().Read(string(path))
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("Unexpectedly couldn't find %s~%s", path, key)
	}

	// Joy of casting in go
	if _, ok := secret.Data["data"]; ok {
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
	}
	return nil, fmt.Errorf("Couldn't find %s~%s", path, key)
}
