package modifier

type valuesTextModifier struct {
}

func (this valuesTextModifier) modify(input []byte) ([]byte, error) {
	list, err := textToKvlist(input)
	if err != nil {
		return nil, err
	}
	return this.modifyKvlist(list)
}

func (_ valuesTextModifier) modifyKvlist(input Kvlist) ([]byte, error) {
	out := []byte{}
	for _, value := range input {
		out = append(out, value.Value...)
	}
	return out, nil
}
