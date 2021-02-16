package substitution

type ValueSource interface {
	GetValue(path []byte, key []byte) (*[]byte, error)
}
