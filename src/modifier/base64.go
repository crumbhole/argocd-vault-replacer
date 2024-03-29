package modifier

import (
	"encoding/base64"
)

type base64Modifier struct{}

func (base64Modifier) modify(input []byte) ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(input)), nil
}
