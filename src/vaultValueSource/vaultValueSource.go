package vaultValueSource

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"log"
	"reflect"
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
	}
	println(m.client.Token())
}

func (m VaultValueSource) GetValue(path []byte, key []byte) *[]byte {
	m.initClient()
	println("GetValue")
	secret, err := m.client.Logical().Read(string(path))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	for i, found := range secret.Data {
		fmt.Printf("Found[%s] %v %v\n", i, found, reflect.TypeOf(found))
	}

	// Joy of casting in go
	switch data := secret.Data["data"].(type) {
	case map[string]interface{}:
		if value, found := data[string(key)]; found {
			fmt.Printf("Got %v %v\n", value, reflect.TypeOf(value))
			switch dataval := value.(type) {
			case string:
				datavalbyte := []byte(dataval)
				return &datavalbyte
			}
		}
	}
	return nil
}
