package main

import (
	"fmt"

	"github.com/blvckbill/lens.git/search"
)

func main() {
	indexer := search.Indexer("doc1.txt", "doc2.txt", "doc3.txt")
	for word, docs := range indexer {
		fmt.Printf("%s → %v\n", word, docs)
	}
}
