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
	DocFreq int
	DocIds  []int
}

type Token struct {
	Word  string
	DocId int
}

type Index struct {
	Docs     []string
	DocToID  map[string]int
	IdToDoc  map[int]string
	Postings map[string]Postings
	NextID   int
}

func NewIndex() *Index {
	return &Index{
		Docs:     make([]string, 0),
		DocToID:  make(map[string]int),
		IdToDoc:  make(map[int]string),
		Postings: make(map[string]Postings),
		NextID:   1,
	}
}

func (idx *Index) BuildFromDir(path string) map[string]Postings {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("Error occurred while trying to locate %s, path not found", path)
	}
	for i := range len(files) {
		idx.Docs = append(idx.Docs, files[i].Name())
	}

	tokens := idx.Tokenize(path, idx.Docs...)
	sortedTokens := idx.SortTokens(tokens)
	postingsList := idx.CompilePostings(sortedTokens)

	return postingsList
}

func (idx *Index) AddDocument(filename string) int {
	_, ok := idx.DocToID[filename]
	if !ok {
		idx.DocToID[filename] = idx.NextID
		idx.IdToDoc[idx.NextID] = filename
		idx.NextID++
	}
	return idx.DocToID[filename]
}

func (idx *Index) Tokenize(path string, documents ...string) []Token {
	var tokens []Token
	for _, document := range documents {
		file := openDocument(path, document)
		docId := idx.AddDocument(file.Name())

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

			tkn := Token{
				Word:  word,
				DocId: docId,
			}

			tokens = append(tokens, tkn)
		}
		file.Close()
	}
	return tokens
}

func (idx *Index) SortTokens(tokens []Token) []Token {
	slices.SortFunc(tokens, func(a, b Token) int {
		if n := strings.Compare(a.Word, b.Word); n != 0 {
			return n
		}
		return cmp.Compare(a.DocId, b.DocId)
	})
	return tokens
}

func (idx *Index) CompilePostings(sortedTokens []Token) map[string]Postings {
	for i := range len(sortedTokens) {
		word := sortedTokens[i].Word
		docId := sortedTokens[i].DocId

		v, ok := idx.Postings[word]

		if !ok {
			idx.Postings[word] = Postings{
				DocFreq: 1,
				DocIds:  []int{docId},
			}
		} else {
			if v.DocIds[len(v.DocIds)-1] == docId {
				continue
			}
			idx.Postings[word] = Postings{
				DocFreq: v.DocFreq + 1,
				DocIds:  append(v.DocIds, docId),
			}
		}
	}
	return idx.Postings
}

func openDocument(path, document string) *os.File {
	file, err := os.OpenInRoot(path, document)
	if err != nil {
		log.Fatalf("Error reading from %s", document)
	}

	return file
}
