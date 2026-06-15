package main

import (
	"fmt"

	"github.com/blvckbill/lens.git/search"
)

func main() {
	indexer := search.Indexer("notes.txt")
	for word, docs := range indexer {
		fmt.Printf("%s → %v\n", word, docs)
	}
}
