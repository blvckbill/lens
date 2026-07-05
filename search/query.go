package search

import (
	"cmp"
	"fmt"
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
	parsedWords := Parser(queryString)
	scores := idx.BuildScore(parsedWords)
	sortedScores := idx.SortScores(scores)
	documents := idx.ResolveDocument(sortedScores)
	for _, s := range sortedScores {
		fmt.Printf("doc: %s score: %.4f\n", idx.IdToDoc[s.DocId], s.Score)
	}

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
// func (idx *Index) PostingsLookup(query_list []string) [][]int {
// 	var postings [][]int
// 	for _, w := range query_list {
// 		v, ok := idx.Postings[w]
// 		if !ok {
// 			log.Printf("No occurence of %s found in documents", w)
// 			break
// 		}
// 		postings = append(postings, v.DocIds)
// 	}
// 	return postings
// }

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
* Builds a list of scores for each document based on the query terms.
 */
func (idx *Index) BuildScore(query_list []string) []ScoreMap {
	NumberOfDocs := len(idx.Docs)
	scores := make(map[int]float64)
	var scoreMap []ScoreMap
	for _, term := range query_list {
		v, ok := idx.Postings[term]
		if !ok {
			log.Printf("No occurence of %s found in documents", term)
			break
		}
		// for the terms in query list, calculate the score of the documents they can be found using N and docfreq
		docFreq := v.DocFreq
		for _, dt := range v.DocTerm {
			scores[dt.DocId] += CalculateScore(dt.TermFreq, NumberOfDocs, docFreq)
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
func (idx *Index) ResolveDocument(doc_list []ScoreMap) []string {
	var documents []string
	for _, scores := range doc_list {
		documents = append(documents, idx.IdToDoc[scores.DocId])
	}
	return documents
}
