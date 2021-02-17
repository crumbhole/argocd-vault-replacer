package substitution

import (
	"fmt"
	"strings"
)

type pathKeyTuple struct {
	path string
	key  string
}

type mockValueSource struct {
	values map[pathKeyTuple][]byte
}

func (m mockValueSource) GetValue(path []byte, key []byte) (*[]byte, error) {
	var pk = pathKeyTuple{strings.TrimSuffix(string(path), `/`), string(key)}
	if val, ok := m.values[pk]; ok {
		return &val, nil
	}
	return nil, fmt.Errorf("Couldn't find %s ! %s", path, key)
}
