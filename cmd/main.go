package main

import (
	"fmt"

	"github.com/blvckbill/lens.git/search"
)

func main() {
	r := search.NewRegistry()

	indexer := r.InvertedIndex("./docs")
	for word, docs := range indexer {
		fmt.Printf("%s → %v\n", word, docs)
	}
}
