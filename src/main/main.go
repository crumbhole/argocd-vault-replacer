package main

import (
	"bytes"
	//"flag"
	"fmt"
	"io/ioutil"
	"log"
	//	vault "github.com/hashicorp/vault/api"
	"os"
	"path/filepath"
	"substitution"
	"vaultValueSource"
)

func updateFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	fileinfo, _ := os.Stat(path)
	origcontents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("File reading error", err)
		return err
	}
	vs := vaultValueSource.VaultValueSource{}
	modifycontents := substitution.Substitute(origcontents, vs)
	if err != nil {
		fmt.Println("Substitution error", err)
		return err
	}
	if !bytes.Equal(modifycontents, origcontents) {
		err = ioutil.WriteFile(path, modifycontents, fileinfo.Mode())
		if err != nil {
			fmt.Println("Couldn't modify file", err)
			return err
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
