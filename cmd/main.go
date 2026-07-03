package main

import (
	"fmt"

	"github.com/blvckbill/lens.git/search"
)

func main() {
	idx := search.NewIndex()

	indexer := idx.BuildFromDir("./docs")
	for word, docs := range indexer {
		fmt.Printf("%s → %v\n", word, docs)
	}
}
