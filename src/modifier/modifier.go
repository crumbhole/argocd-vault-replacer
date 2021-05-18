package modifier

type Kv struct {
	Key   []byte
	Value []byte
}

type Kvlist []Kv

type modifier interface {
	modify([]byte) ([]byte, error)
}

type kvmodifier interface {
	modifier
	modifyKvlist(Kvlist) ([]byte, error)
}
