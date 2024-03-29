package modifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type htaccessModifier struct{}

func (htaccessModifier) modify(inputJSON []byte) ([]byte, error) {
	input := make([]map[string]string, 0)
	err := json.Unmarshal(inputJSON, &input)
	if err != nil {
		inputsingle := make(map[string]string, 0)
		err = json.Unmarshal(inputJSON, &inputsingle)
		if err != nil {
			return inputJSON, err
		}
		input = append(input, inputsingle)
	}
	passwords := make(map[string]string, len(input))
	for _, kv := range input {
		user, ok := kv[`user`]
		if !ok {
			return nil, errors.New(`No key called user in input json`)
		}
		password, ok := kv[`password`]
		if !ok {
			return nil, errors.New(`No key called password in input json`)
		}
		hashedpw, err := htaccessBcrypt(password)
		if err != nil {
			return nil, err
		}
		passwords[user] = string(hashedpw)
	}
	return htaccessEncode(passwords), nil
}

func htaccessEncode(in map[string]string) (out []byte) {
	out = []byte{}
	for name, pw := range in {
		out = append(out, []byte(fmt.Sprintf("%s:%s\n", name, pw))...)
	}
	return out
}

func htaccessBcrypt(password string) (hash []byte, err error) {
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
