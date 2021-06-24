package bwValueSource

import (
	"errors"
	"os"
	"strings"

	bwwrap "github.com/crumbhole/bitwardenwrapper"
)

const (
	envCheck = "BW_SESSION"
)

type BitwardenValueSource struct{}

func (_ BitwardenValueSource) getItemSplitPath(path string) (*bwwrap.BwItem, error) {
	pathParts := strings.Split(string(path), `/`)
	keyUsed := pathParts[len(pathParts)-1]
	pathUsed := strings.Join(pathParts[:len(pathParts)-1], `/`)
	return bwwrap.GetItemFromFolder(keyUsed, pathUsed)
}

func (m BitwardenValueSource) GetValue(path []byte, key []byte) (*[]byte, error) {
	if _, present := os.LookupEnv(envCheck); !present {
		return nil, errors.New("Bitwarden session key not present")
	}
	switch string(key) {
	default:
		item, err := bwwrap.GetItemFromFolder(string(key), string(path))
		if err != nil {
			return nil, err
		}
		value := []byte(item.Notes)
		return &value, nil
	case `username`:
		item, err := m.getItemSplitPath(string(path))
		if err != nil {
			return nil, err
		}
		value := []byte(item.Login.Username)
		return &value, nil
	case `password`:
		item, err := m.getItemSplitPath(string(path))
		if err != nil {
			return nil, err
		}
		value := []byte(item.Login.Password)
		return &value, nil
	}
}
