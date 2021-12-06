package substitution

// ValueSource is the interface defining a single call
// GetValue - takes a path and key to a value and returns that value, or null and an error explaining why not
type ValueSource interface {
	GetValue(path []byte, key []byte) (*[]byte, error)
}
