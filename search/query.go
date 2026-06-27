package search

import (
	"log"
	"strings"
)

func Parser(query_string string) []string {
	return strings.Fields(query_string)
}

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

func getIntersection(postings [][]int) []int {
	var intersection []int
	list1 := postings[0]
	list2 := postings[1]

	ptr1 := 0
	ptr2 := 0

	if len(list2) < len(list1) {
		list1, list2 = list2, list1
	}

	for ptr1 < len(list1) {
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

func (idx *Index) resolveDocument(doc_list []int) []string {
	var documents []string
	for _, docId := range doc_list {
		documents = append(documents, idx.IdToDoc[docId])
	}
	return documents
}
