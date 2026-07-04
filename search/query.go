package search

import (
	"log"
	"strings"
)

/**
* Run query by parsing the input words into a list then get their documentIds and
* return the documents where they can be found after intersecting their docIds.
 */
func (idx *Index) Query(queryString string) []string {
	parsedWords := Parser(queryString)
	documentIds := idx.PostingsLookup(parsedWords)
	intersection := IntersectManyLists(documentIds)
	documents := idx.ResolveDocument(intersection)

	return documents
}

/**
*	Receive a query string, parse it into a list of words, and return the list of words.
 */
func Parser(query_string string) []string {
	return strings.Fields(query_string)
}

/**
* Loops through the list of words in the query, and for each word, looks up the postings list in the index.
* If the word is not found in the index, it logs a message and breaks out of the loop.
* Returns a list of doc ids for each word in the query.
 */
func (idx *Index) PostingsLookup(query_list []string) [][]int {
	var postings [][]int
	for _, w := range query_list {
		v, ok := idx.Postings[w]
		if !ok {
			log.Printf("No occurence of %s found in documents", w)
			break
		}
		postings = append(postings, v.DocIds)
	}
	return postings
}

/**
* Lopps throught the two list of docIds and returns the intersection of the two lists.
* The intersection is the list of docIds where you can find the two words in the query.
 */
func IntersectTwoLists(list1, list2 []int) []int {
	var intersection []int

	ptr1 := 0
	ptr2 := 0

	for ptr1 < len(list1) && ptr2 < len(list2) {
		if list1[ptr1] == list2[ptr2] {
			intersection = append(intersection, list1[ptr1])
			ptr1++
			ptr2++
		} else if list1[ptr1] < list2[ptr2] {
			ptr1++
		} else {
			ptr2++
		}
	}
	return intersection
}

/**
* if the query contains more than two words, it find the intersection of list of docids for each word in the query.
* It returns the list of docIds where you can find all the words in the query.
 */
func IntersectManyLists(postings [][]int) []int {
	if len(postings) == 0 {
		return nil
	}

	if len(postings) == 1 {
		return postings[0]
	}
	current := postings[0]
	for i := 1; i < len(postings); i++ {
		current = IntersectTwoLists(current, postings[i])
	}
	return current
}

/**
* Takes a list of document IDs and returns the corresponding document names.
 */
func (idx *Index) ResolveDocument(doc_list []int) []string {
	var documents []string
	for _, docId := range doc_list {
		documents = append(documents, idx.IdToDoc[docId])
	}
	return documents
}
