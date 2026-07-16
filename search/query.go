package search

import (
	"cmp"
	"log"
	"math"
	"slices"
	"strings"
)

type ScoreMap struct {
	Score float64
	DocId int
}

/**
* Run query by parsing the input words into a list then get their documentIds and
* return the documents where they can be found after intersecting their docIds.
 */
func (idx *Index) Query(queryString string) []string {
	parsedTerms := Parser(queryString)

	if len(parsedTerms) != 3 {
		log.Printf("invalid query %q: expected '<term> <operator> <term>'", queryString)
		return nil
	}

	left, ok := idx.Lookup(parsedTerms[0])
	if !ok {
		log.Printf("term %q not found", parsedTerms[0])
		return nil
	}

	right, ok := idx.Lookup(parsedTerms[2])
	if !ok {
		log.Printf("term %q not found", parsedTerms[2])
		return nil
	}

	var docIDs []int

	switch strings.ToUpper(parsedTerms[1]) {
	case "AND":
		docIDs = Intersect(left.DocIDs(), right.DocIDs())

	case "OR":
		docIDs = Union(left.DocIDs(), right.DocIDs())

	case "NOT":
		docIDs = Difference(left.DocIDs(), right.DocIDs())

	default:
		log.Printf("invalid operator %q", parsedTerms[1])
		return nil
	}

	return idx.ResolveDocuments(docIDs)
}

/**
*	Receive a query string, parse it into a list of words, and return the list of words.
 */
func Parser(query_string string) []string {
	return strings.Fields(query_string)
}

// Lookup returns the index entry for a term.
func (idx *Index) Lookup(term string) (TermInfo, bool) {
	info, ok := idx.InvertedIndex[term]
	return info, ok
}

// DocIDs returns the document IDs that contain the term.
func (t *TermInfo) DocIDs() []int {
	ids := make([]int, 0, len(t.Postings))

	for _, posting := range t.Postings {
		ids = append(ids, posting.DocID)
	}

	return ids
}

/**
* Lopps throught the two list of docIds and returns the intersection of the two lists.
* The intersection is the list of docIds where you can find the two words in the query.
 */
func Intersect(list1, list2 []int) []int {
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
func IntersectAll(postings [][]int) []int {
	if len(postings) == 0 {
		return nil
	}

	if len(postings) == 1 {
		return postings[0]
	}
	current := postings[0]
	for i := 1; i < len(postings); i++ {
		current = Intersect(current, postings[i])
	}
	return current
}

func Union(list1, list2 []int) []int {
	var union []int
	ptr1 := 0
	ptr2 := 0

	for ptr1 < len(list1) && ptr2 < len(list2) {
		if list1[ptr1] == list2[ptr2] {
			union = append(union, list1[ptr1])
			ptr1++
			ptr2++
		} else if list1[ptr1] < list2[ptr2] {
			union = append(union, list1[ptr1])
			ptr1++
		} else {
			union = append(union, list2[ptr2])
			ptr2++
		}
	}
	if ptr1 < len(list1) {
		for ptr1 < len(list1) {
			union = append(union, list1[ptr1])
			ptr1++
		}
	}

	if ptr2 < len(list2) {
		for ptr2 < len(list2) {
			union = append(union, list2[ptr2])
			ptr2++
		}
	}
	return union
}

func Difference(list1, list2 []int) []int {
	var difference []int
	ptr1 := 0
	ptr2 := 0

	for ptr1 < len(list1) && ptr2 < len(list2) {
		if list1[ptr1] == list2[ptr2] {
			ptr1++
			ptr2++
		} else if list1[ptr1] < list2[ptr2] {
			difference = append(difference, list1[ptr1])
			ptr1++
		} else {
			ptr2++
		}
	}
	if ptr1 < len(list1) {
		for ptr1 < len(list1) {
			difference = append(difference, list1[ptr1])
			ptr1++
		}
	}
	return difference
}

/**
* Builds a list of scores for each document based on the query terms.
 */
func (idx *Index) BuildScore(query_list []string) []ScoreMap {
	NumberOfDocs := len(idx.Documents)
	scores := make(map[int]float64)
	var scoreMap []ScoreMap
	for _, term := range query_list {
		v, ok := idx.InvertedIndex[term]
		if !ok {
			log.Printf("No occurence of %s found in documents", term)
			break
		}
		// for the terms in query list, calculate the score of the documents they can be found using N and docfreq
		docFreq := v.DocumentFrequency
		for _, dt := range v.Postings {
			scores[dt.DocID] += CalculateScore(len(dt.Positions), NumberOfDocs, docFreq)
		}
	}
	for docId, score := range scores {
		s := ScoreMap{
			Score: score,
			DocId: docId,
		}
		scoreMap = append(scoreMap, s)
	}
	return scoreMap
}

/**
* Calculates the score of a document based on the term frequency, number of documents, and document frequency.
* The score is calculated using the formula: score = tf * log10(N / df)
 */
func CalculateScore(tf, N, df int) float64 {
	invertedDocFreq := float64(N) / float64(df)
	return float64(tf) * math.Log10(invertedDocFreq)
}

/**
* Sorts the list of scores in descending order based on the score value.
 */
func (idx *Index) SortScores(scores []ScoreMap) []ScoreMap {
	slices.SortFunc(scores, func(a, b ScoreMap) int {
		if n := cmp.Compare(b.Score, a.Score); n != 0 {
			return n
		}
		return cmp.Compare(a.DocId, b.DocId)
	})
	return scores
}

/**
* Takes a list of document IDs and returns the corresponding document names.
 */
func (idx *Index) ResolveDocuments(doc_list []int) []string {
	var documents []string
	for _, docId := range doc_list {
		documents = append(documents, idx.IDToDocument[docId])
	}
	return documents
}
