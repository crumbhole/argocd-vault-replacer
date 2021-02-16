package main

import (
	"fmt"
	"github.com/joibel/vault-replacer/src/substitution"
	"github.com/joibel/vault-replacer/src/vaultValueSource"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func scanFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	fileRegexp := regexp.MustCompile(`\.ya?ml$`)
	if fileRegexp.MatchString(path) {
		origcontents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		modifycontents, err := substitution.Substitute(origcontents, vaultValueSource.VaultValueSource{})
		if err != nil {
			return err
		}
		fmt.Printf("---\n%s\n", modifycontents)
	}
	return nil
}

func scanDir(path string) error {
	return filepath.Walk(path, scanFile)
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = scanDir(dir)
	if err != nil {
		log.Fatal(err)
	}
}
