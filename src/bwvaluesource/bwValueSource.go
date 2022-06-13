package bwvaluesource

import (
	"errors"
	"os"
	"strings"

	bwwrap "github.com/crumbhole/bitwardenwrapper"
)

const (
	envCheck   = "BW_SESSION"
	argoPrefix = "ARGOCD_ENV_"
)

// BwSession returns true of BW_SESSION or ARGOCD_ENV_BW_SESSION are set
// If ARGOCD_ENV_BWSESSION is set the value is copied to BW_SESSION
func BwSession() bool {
	val, got := os.LookupEnv(argoPrefix + envCheck)
	if !got {
		val, got = os.LookupEnv(envCheck)
		if !got {
			return false
		}
		return true
	}
	os.Setenv(envCheck, val)
	return true
}

// BitwardenValueSource is a value source getting values from bitwarden
type BitwardenValueSource struct{}

func (BitwardenValueSource) getItemSplitPath(path string) (*bwwrap.BwItem, error) {
	pathParts := strings.Split(string(path), `/`)
	keyUsed := pathParts[len(pathParts)-1]
	pathUsed := strings.Join(pathParts[:len(pathParts)-1], `/`)
	return bwwrap.GetItemFromFolder(keyUsed, pathUsed)
}

// GetValue returns a value from a path+key in bitwarden or null if it doesn't exist
func (m BitwardenValueSource) GetValue(path []byte, key []byte) (*[]byte, error) {
	if !BwSession() {
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
