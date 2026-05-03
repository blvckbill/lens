package search

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Search_file(arg string) {
	file, err := os.Open("notes.txt")
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineNumber := 1
	for scanner.Scan() {
		lineText := scanner.Text()
		if strings.Contains(lineText, arg) {
			fmt.Printf("Found '%s' on line %d\n", arg, lineNumber)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error scanning file: ", err)
	}
}
