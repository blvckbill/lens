package search

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func Search_file(arg string) {
	file, err := os.ReadFile("notes.txt")
	if err != nil {
		log.Fatal("An error occured while opening this file")
	}

	fileContent := string(file)

	if strings.Contains(fileContent, arg) {
		fmt.Println("Word found in file")
	} else {
		fmt.Println("Word not found in file")
	}
}
