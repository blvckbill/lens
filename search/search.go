package search

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func Search_file(document string, word string) bool {
	file, err := os.Open(document)
	if err != nil {
		log.Fatal("Error reading from file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		read_text := strings.ToLower(scanner.Text())
		if strings.Contains(read_text, word) {
			return true
		}
	}
	return false
}
