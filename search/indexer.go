package search

import (
	"bufio"
	"cmp"
	"log"
	"os"
	"slices"
	"strings"
	"unicode"
)

type Postings struct {
	docFreq int
	docIds  []int
}

type Registry struct {
	docs     []string
	docToID  map[string]int
	idTodoc  map[int]string
	postings map[string]Postings
	nextID   int
}

type WordIndex struct {
	word  string
	docId int
}

func NewRegistry() *Registry {
	return &Registry{
		docs:     make([]string, 0),
		docToID:  make(map[string]int),
		idTodoc:  make(map[int]string),
		postings: make(map[string]Postings),
		nextID:   1,
	}
}

func (r *Registry) InvertedIndex(path string) map[string]Postings {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Error occured while trying to locate %s, path not found", path)
	}
	for i := range len(files) {
		r.docs = append(r.docs, files[i].Name())
	}

	wordIndex := r.PairWordsAndDocId(path, r.docs...)
	sortedPairs := r.SortPairsByWord(wordIndex)
	postingsList := r.TermDictionary(sortedPairs)

	return postingsList
}

func (r *Registry) AddDocument(filename string) int {
	_, ok := r.docToID[filename]
	if !ok {
		r.docToID[filename] = r.nextID
		r.idTodoc[r.nextID] = filename
		r.nextID++
	}
	return r.docToID[filename]
}

func (r *Registry) PairWordsAndDocId(path string, documents ...string) []WordIndex {
	var pairs []WordIndex
	for _, document := range documents {
		file := openDocument(path, document)
		// map doc to id and return the id
		docId := r.AddDocument(file.Name())

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			var b strings.Builder
			text := scanner.Text()
			for _, t := range text {
				if unicode.IsPunct(t) {
					continue
				} else {
					b.WriteString(string(t))
				}
			}
			word := strings.ToLower(b.String())

			if word == "" {
				continue
			}

			// map word to doc id
			w := WordIndex{
				word:  word,
				docId: docId,
			}

			// append the pair (word, docid) to slice pairs
			pairs = append(pairs, w)
		}
		file.Close()
	}
	// return list of all pairs
	return pairs
}

func (r *Registry) SortPairsByWord(pairs []WordIndex) []WordIndex {
	slices.SortFunc(pairs, func(a, b WordIndex) int {
		if n := strings.Compare(a.word, b.word); n != 0 {
			return n
		}
		return cmp.Compare(a.docId, b.docId)
	})
	return pairs
}

func (r *Registry) TermDictionary(sorted_pairs []WordIndex) map[string]Postings {
	for i := range len(sorted_pairs) {
		word := sorted_pairs[i].word
		docId := sorted_pairs[i].docId

		v, ok := r.postings[word]

		if !ok {
			r.postings[word] = Postings{
				docFreq: 1,
				docIds:  []int{docId},
			}
		} else {
			if v.docIds[len(v.docIds)-1] == docId {
				continue
			}
			r.postings[word] = Postings{
				docFreq: v.docFreq + 1,
				docIds:  append(v.docIds, docId),
			}
		}
	}
	return r.postings
}

func openDocument(path, document string) *os.File {
	file, err := os.OpenInRoot(path, document)
	if err != nil {
		log.Fatalf("Error reading from %s", document)
	}

	return file
}
