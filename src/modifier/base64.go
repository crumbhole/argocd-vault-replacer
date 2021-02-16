package modifier

import (
	"encoding/base64"
)

type base64Modifier struct{}

func (_ base64Modifier) modify(input []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(input))
}
