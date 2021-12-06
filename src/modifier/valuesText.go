package modifier

type valuesTextModifier struct {
}

func (mod valuesTextModifier) modify(input []byte) ([]byte, error) {
	list, err := textToKvlist(input)
	if err != nil {
		return nil, err
	}
	return mod.modifyKvlist(list)
}

func (valuesTextModifier) modifyKvlist(input Kvlist) ([]byte, error) {
	out := []byte{}
	for _, value := range input {
		out = append(out, value.Value...)
	}
	return out, nil
}
