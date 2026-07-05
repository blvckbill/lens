package main

import (
	"fmt"
	"path/filepath"

	"github.com/blvckbill/lens.git/search"
)

func main() {
	idx := search.NewIndex()

	indexer := idx.BuildFromDir("./docs")
	for word, docs := range indexer {
		fmt.Printf("%s → %v\n", word, docs)
	}
	result := idx.Query("dog fast runs")
	fmt.Println(result)
	for _, filePath := range result {
		filename := filepath.Base(filePath)
		fmt.Printf("The words can be found in %s", filename)
	}
	for word, posting := range idx.Postings {
		for _, dt := range posting.DocTerm {
			if dt.TermFreq > 1 {
				fmt.Printf("%s → doc %d → freq %d\n", word, dt.DocId, dt.TermFreq)
			}
		}
	}

}
