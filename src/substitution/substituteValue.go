package substitution

import (
	"fmt"
	"log"
	"regexp"
)

func substituteValue(input []byte, source ValueSource) []byte {
	reOuter := regexp.MustCompile(`^<\s*value:\s*([^\!]*[^\s])\s*(\!\s*[^\|]+)?\s*(\|.*)?\s*>$`)
	pathFound := reOuter.FindSubmatch(input)
	for i, found := range pathFound{
		fmt.Printf("Found[%d] %s\n", i, found)
	}
	if pathFound != nil {
		if len(pathFound[2]) > 0 {
			fmt.Printf("Key found %d\n", len(pathFound))
			reKey := regexp.MustCompile(`^!\s*(.*?)\s*$`)
			fmt.Printf("Not found //%v//\n", pathFound[2])
			keyFound := reKey.FindSubmatch(pathFound[2])
			if keyFound == nil {
				log.Fatal("Key regex failure")
				fmt.Printf("Key regex failure\n")
				return input
			}
			fmt.Printf("Success %v\n", keyFound[1])
			value := source.GetValue(pathFound[1], keyFound[1])
			if value == nil {
				log.Fatal("Key lookup failure")
				fmt.Printf("Key lookup failure\n")
				return input
			}
			return *value
		}
		fmt.Printf("No key found %d\n", len(pathFound))
		return pathFound[1]
	}
	println("Not match")
	return input
}
