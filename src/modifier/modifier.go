package modifier

// Kv is a Key Value Pair
type Kv struct {
	Key   []byte
	Value []byte
}

// Kvlist is a list of key values
type Kvlist []Kv

type modifier interface {
	modify([]byte) ([]byte, error)
}

type kvmodifier interface {
	modifier
	modifyKvlist(Kvlist) ([]byte, error)
}
