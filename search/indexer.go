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

type TermFreq struct {
	TermFreq int
	DocId int
}

type Postings struct {
	DocFreq int
	DocTerm  []TermFreq
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

/*
*
Create a new index, which is a data structure that maps words to the documents they appear in.
*/
func NewIndex() *Index {
	return &Index{
		Docs:     make([]string, 0),
		DocToID:  make(map[string]int),
		IdToDoc:  make(map[int]string),
		Postings: make(map[string]Postings),
		NextID:   1,
	}
}

/**	Receive the path to the directory containing the documents,
*	and build the index from the documents in that directory.
 */
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

/*
* map document to an integer ID, if the document is already mapped,
return the existing ID, otherwise assign a new ID and return it.
*/
func (idx *Index) AddDocument(filename string) int {
	_, ok := idx.DocToID[filename]
	if !ok {
		idx.DocToID[filename] = idx.NextID
		idx.IdToDoc[idx.NextID] = filename
		idx.NextID++
	}
	return idx.DocToID[filename]
}

/*
*
loop through the documents, open each document, read its content, and tokenize the words in the document.
*/
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

/*
*
Loop through the tokenized documents, which maps word to document ID
and sorts it so same words are grouped together.
*/
func (idx *Index) SortTokens(tokens []Token) []Token {
	slices.SortFunc(tokens, func(a, b Token) int {
		if n := strings.Compare(a.Word, b.Word); n != 0 {
			return n
		}
		return cmp.Compare(a.DocId, b.DocId)
	})
	return tokens
}

/*
*
Loop through the sorted tokens, and build the postings list,
which maps each word to the number of times it appears globally and list of its document IDs.
*/
func (idx *Index) CompilePostings(sortedTokens []Token) map[string]Postings {
	for i := range len(sortedTokens) {
		word := sortedTokens[i].Word
		docId := sortedTokens[i].DocId

		v, ok := idx.Postings[word]

		tf := TermFreq{
			TermFreq: 1,
			DocId: docId,
		}

		if !ok {
			idx.Postings[word] = Postings{
				DocFreq: 1,
				DocTerm:  []TermFreq{tf},
			}
		} else {
			// if last appended term frequency struct has same doc id, increment term frequency
			if v.DocTerm[len(v.DocTerm)-1].DocId == docId {
				v.DocTerm[len(v.DocTerm)-1].TermFreq++
				idx.Postings[word] = v
				continue
			}
			// if the word is found in the dictionary, append the term freq and doc id and then increment the doc freq
			idx.Postings[word] = Postings{
				DocFreq: v.DocFreq + 1,
				DocTerm:  append(v.DocTerm, tf),
			}
		}
	}
	return idx.Postings
}

/*
*
Open the document in the specified path and return a file pointer to it.
*/
func openDocument(path, document string) *os.File {
	file, err := os.OpenInRoot(path, document)
	if err != nil {
		log.Fatalf("Error reading from %s", document)
	}

	return file
}
