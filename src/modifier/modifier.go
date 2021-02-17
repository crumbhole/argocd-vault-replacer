package modifier

type modifier interface {
	modify([]byte) ([]byte, error)
}
