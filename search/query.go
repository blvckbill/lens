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

func (idx *Index) ResolveDocument(doc_list []int) []string {
	var documents []string
	for _, docId := range doc_list {
		documents = append(documents, idx.IdToDoc[docId])
	}
	return documents
}
