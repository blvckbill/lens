package main

import (
	"fmt"
	//"path/filepath"

	"github.com/blvckbill/lens.git/search"
)

func main() {
	idx := search.NewIndex()

	indexer := idx.BuildFromDir("./docs")
	for word, docs := range indexer {
		fmt.Printf("%s → %v\n", word, docs)
	}
	// result := idx.Query("dog fast runs")
	// fmt.Println(result)
	// for _, filePath := range result {
	// 	filename := filepath.Base(filePath)
	// 	fmt.Printf("The words can be found in %s", filename)
	// }
	// for word, terminfo := range idx.InvertedIndex {
	// 	for _, dt := range terminfo.Postings {
	// 		if len(dt.Positions) > 1 {
	// 			fmt.Printf("%s → doc %d → freq %d\n", word, dt.DocID, len(dt.Positions))
	// 		}
	// 	}
	// }

}
