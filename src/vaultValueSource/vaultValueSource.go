package vaultValueSource

type pathKeyTuple struct {
	path string
	key  string
}

type VaultValueSource struct {
	values map[pathKeyTuple][]byte
}

func (m VaultValueSource) GetValue(path []byte, key []byte) *[]byte {
	var pk = pathKeyTuple{string(path), string(key)}
	if val, ok := m.values[pk]; ok {
		return &val
	}
	return nil
}
