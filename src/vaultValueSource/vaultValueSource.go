package vaultValueSource

import (
	vault "github.com/hashicorp/vault/api"
	"log"
)

type VaultValueSource struct {
	client *vault.Client
}

func (m *VaultValueSource) initClient() {
	if m.client == nil {
		client, err := vault.NewClient(nil)
		if err != nil {
			log.Fatal(err)
		}
		m.client = client
		m.tryKubernetesAuth()
	}
}

func (m VaultValueSource) GetValue(path []byte, key []byte) *[]byte {
	m.initClient()
	secret, err := m.client.Logical().Read(string(path))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Joy of casting in go
	switch data := secret.Data["data"].(type) {
	case map[string]interface{}:
		if value, found := data[string(key)]; found {
			switch dataval := value.(type) {
			case string:
				datavalbyte := []byte(dataval)
				return &datavalbyte
			}
		}
	}
	return nil
}
