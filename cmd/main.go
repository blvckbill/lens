package main

import (
	"fmt"

	"github.com/blvckbill/lens.git/search"
)

func main() {
	found := search.Search_file("notes.txt", "large")
	if found {
		fmt.Println("found a match")
	} else {
		fmt.Println("No found matches")
	}

}
