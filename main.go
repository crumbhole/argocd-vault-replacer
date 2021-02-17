package main

import (
	"fmt"
	"github.com/joibel/argocd-vault-replacer/src/substitution"
	"github.com/joibel/argocd-vault-replacer/src/vaultValueSource"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type scanner struct {
	source substitution.ValueSource
}

func (s *scanner) scanFile(path string, info os.FileInfo, err error) error {
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
		modifycontents, err := substitution.Substitute(origcontents, s.source)
		if err != nil {
			return err
		}
		fmt.Printf("---\n%s\n", modifycontents)
	}
	return nil
}

func (s *scanner) scanDir(path string) error {
	return filepath.Walk(path, s.scanFile)
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	s := scanner{source: vaultValueSource.VaultValueSource{}}
	err = s.scanDir(dir)
	if err != nil {
		log.Fatal(err)
	}
}
