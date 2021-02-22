package modifier

type Kv struct {
	Key   []byte
	Value []byte
}

type Kvlist []Kv

type kvmodifier interface {
	modify(Kvlist) ([]byte, error)
}

type modifier interface {
	modify([]byte) ([]byte, error)
}
