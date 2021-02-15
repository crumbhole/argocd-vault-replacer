package main

import (
	"bytes"
	//"flag"
	"fmt"
	"io/ioutil"
	"log"
	//	vault "github.com/hashicorp/vault/api"
	"github.com/joibel/vault-replacer/src/substitution"
	"github.com/joibel/vault-replacer/src/vaultValueSource"
	"os"
	"path/filepath"
	"regexp"
)

func updateFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	if info.IsDir() {
		return nil
	}
	fileRegexp := regexp.MustCompile(`\.ya?ml$`)
	if fileRegexp.MatchString(path) {
		origcontents, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("GGFile reading error", err)
			return err
		}
		modifycontents := substitution.Substitute(origcontents, vaultValueSource.VaultValueSource{})
		if err != nil {
			fmt.Println("Substitution error", err)
			return err
		}
		if !bytes.Equal(modifycontents, origcontents) {
			err = ioutil.WriteFile(path, modifycontents, info.Mode())
			if err != nil {
				fmt.Println("Couldn't modify file", err)
				return err
			}
		}
	}
	return nil
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(dir, updateFile)
	if err != nil {
		log.Fatal(err)
	}
}
