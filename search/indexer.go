package search

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"
)

func Indexer(documents ...string) map[string]map[string][]int {
	var indexer = make(map[string][]map[string][]int, 10)
	for _, document := range documents {
		file := openDocument(document)

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)

		lineNumber := 1
		for scanner.Scan() {
			var b strings.Builder
			text := scanner.Text()
			for _, t := range text {
				if unicode.IsPunct(t) {
					continue
				} else {
					b.WriteString(string(t))
				}
			}
			word := strings.ToLower(b.String())
			alreadyExists := false
			var docMap = make(map[string][]int)
			for doc := range indexer[word] {
				if document == doc {
					alreadyExists = true
				}
			}
			if alreadyExists {
				docMap[document] = append(docMap[document], lineNumber)
			} else {
				docMap[document] = append(docMap[document], lineNumber)
				indexer[word] = docMap
			}
		}
		file.Close()
	}
	return indexer
}

func openDocument(document string) *os.File {
	file, err := os.Open(document)
	if err != nil {
		log.Fatalf("Error reading from document %s", document)
	}

	return file
}
