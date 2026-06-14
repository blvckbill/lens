package search

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"
)

func Indexer(documents ...string) map[string]map[string]int {
	var indexer = make(map[string]map[string]int, 10)
	for _, document := range documents {
		file := openDocument(document)

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)

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

			if indexer[word] == nil {
				indexer[word] = make(map[string]int)
			}
			indexer[word][document] += 1

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
