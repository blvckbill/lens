package search

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"
)

func Indexer(documents ...string) map[string][]string {
	var indexer = make(map[string][]string, 10)
	for _, document := range documents {
		file := openDocument(document)

		scanner := bufio.NewScanner(file)

		
		for scanner.Scan() {
			var cleanedText string
			text := scanner.Text()
			for _, t := range text {
				if unicode.IsPunct(t) {
					continue
				} else {
					cleanedText += string(t)
				}
			}
			normalizedLine := strings.ToLower(cleanedText)
			words := strings.SplitSeq(normalizedLine, " ")
			for word := range words {
				alreadyExists := false
				for _, doc := range indexer[word] {
					if doc == document {
						alreadyExists = true
					}
					
				}
				if !alreadyExists {
					indexer[word] = append(indexer[word], document)
				}
			}
		}
		file.Close()
	}
	return indexer
}

func openDocument(document string) *os.File{
	file, err := os.Open(document)
	if err != nil {
		log.Fatalf("Error reading from document %s", document)
	}


	return file
}