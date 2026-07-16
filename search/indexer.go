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

// Token represents a normalized word extracted from a document.
type Token struct {
	Word     string
	DocID    int
	Position int
}

// Posting stores information about a term within a single document.
type Posting struct {
	DocID     int
	Positions []int
}

// TermInfo stores all postings for a term.
type TermInfo struct {
	DocumentFrequency int
	Postings          []Posting
}

// Index represents an in-memory positional inverted index.
type Index struct {
	Documents      []string
	DocumentIDs    map[string]int
	IDToDocument   map[int]string
	InvertedIndex  map[string]TermInfo
	NextDocumentID int
}

// NewIndex creates an empty index.
func NewIndex() *Index {
	return &Index{
		Documents:      make([]string, 0),
		DocumentIDs:    make(map[string]int),
		IDToDocument:   make(map[int]string),
		InvertedIndex:  make(map[string]TermInfo),
		NextDocumentID: 1,
	}
}

// BuildFromDir indexes every document in the supplied directory.
func (idx *Index) BuildFromDir(path string) map[string]TermInfo {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("unable to read directory %q", path)
		return nil
	}

	for _, file := range files {
		idx.Documents = append(idx.Documents, file.Name())
	}

	tokens := idx.Tokenize(path, idx.Documents...)
	sortedTokens := idx.SortTokens(tokens)

	return idx.CompilePostings(sortedTokens)
}

// AddDocument assigns a unique ID to a document.
func (idx *Index) AddDocument(filename string) int {
	if id, exists := idx.DocumentIDs[filename]; exists {
		return id
	}

	id := idx.NextDocumentID

	idx.DocumentIDs[filename] = id
	idx.IDToDocument[id] = filename
	idx.NextDocumentID++

	return id
}

// Tokenize converts each document into normalized tokens.
func (idx *Index) Tokenize(path string, documents ...string) []Token {
	var tokens []Token

	for _, document := range documents {
		file := openDocument(path, document)
		docID := idx.AddDocument(file.Name())

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)

		position := 0

		for scanner.Scan() {
			var builder strings.Builder

			for _, r := range scanner.Text() {
				if unicode.IsPunct(r) {
					continue
				}
				builder.WriteRune(unicode.ToLower(r))
			}

			word := builder.String()
			if word == "" {
				continue
			}

			tokens = append(tokens, Token{
				Word:     word,
				DocID:    docID,
				Position: position,
			})

			position++
		}

		file.Close()
	}

	return tokens
}

// SortTokens sorts tokens by term then by document ID then lastly by position.
func (idx *Index) SortTokens(tokens []Token) []Token {
	slices.SortFunc(tokens, func(a, b Token) int {
		if n := strings.Compare(a.Word, b.Word); n != 0 {
			return n
		}

		if n := cmp.Compare(a.DocID, b.DocID); n != 0 {
			return n
		}

		return cmp.Compare(a.Position, b.Position)
	})

	return tokens
}

// CompilePostings builds the positional inverted index.
func (idx *Index) CompilePostings(tokens []Token) map[string]TermInfo {
	for _, token := range tokens {
		term := token.Word

		entry, exists := idx.InvertedIndex[term]

		if !exists {
			idx.InvertedIndex[term] = TermInfo{
				DocumentFrequency: 1,
				Postings: []Posting{
					{
						DocID:     token.DocID,
						Positions: []int{token.Position},
					},
				},
			}
			continue
		}

		lastPosting := &entry.Postings[len(entry.Postings)-1]

		if lastPosting.DocID == token.DocID {
			lastPosting.Positions = append(lastPosting.Positions, token.Position)
		} else {
			entry.DocumentFrequency++
			entry.Postings = append(entry.Postings, Posting{
				DocID:     token.DocID,
				Positions: []int{token.Position},
			})
		}

		idx.InvertedIndex[term] = entry
	}

	return idx.InvertedIndex
}

// openDocument opens a document from the collection.
func openDocument(path, document string) *os.File {
	file, err := os.OpenInRoot(path, document)
	if err != nil {
		log.Fatalf("error opening %s", document)
	}

	return file
}
