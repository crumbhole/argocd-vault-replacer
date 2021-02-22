package modifier

type valuesTextModifier struct{}

func (_ valuesTextModifier) modify(input Kvlist) ([]byte, error) {
	out := []byte{}
	for _, value := range input {
		out = append(out, value.Value...)
	}
	return out, nil
}
