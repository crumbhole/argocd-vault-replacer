package substitution

import (
	//	"fmt"
	"strings"
)

type pathKeyTuple struct {
	path string
	key  string
}

type mockValueSource struct {
	values map[pathKeyTuple][]byte
}

func (m mockValueSource) GetValue(path []byte, key []byte) *[]byte {
	var pk = pathKeyTuple{strings.TrimSuffix(string(path), `/`), string(key)}
	//	fmt.Printf("Looking up %v\n", pk)
	if val, ok := m.values[pk]; ok {
		//	fmt.Printf("Mock Lookup found %s\n", val)
		return &val
	}
	return nil
}
