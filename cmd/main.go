package main

// local files search engine using pagerank algorithm

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// FileInfo struct to store term frequency, file name, word count, size, last modified time
type FileInfo struct {
	TermFrequency map[string]int
	// FilePath      string
	WordCount    int
	Size         int64
	LastModified int64
}

// Indexer struct to store map of file info, file frequency
type Indexer struct {
	FileInfo map[string]FileInfo
	FileFreq map[string]int
}

type Loc struct {
	Row int
	Col int
}

type TermLoc struct {
	FilePath string
	Lines    []Loc
}

// Indexer constructor
func NewIndexer() *Indexer {
	return &Indexer{
		FileInfo: make(map[string]FileInfo),
		FileFreq: make(map[string]int),
	}
}

func (i *Indexer) IndexFiles(dir string, extensions []string) {
	// index files in directory recursively
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			for _, e := range extensions {
				if ext == e {
					i.AddFile(path)
					break
				}
			}
		}
		return nil
	})
}

// Indexer method to add file to index
func (i *Indexer) AddFile(filePath string) {
	// add file to index
	fileStat, err := os.Stat(filePath)
	if err != nil {
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	fi := NewFileInfo()
	fi.Size = fileStat.Size()
	fi.LastModified = fileStat.ModTime().Unix()

	data := make([]byte, fi.Size)
	_, err = file.Read(data)
	if err != nil {
		return
	}

	content := string(data)
	words := strings.Fields(content)
	fi.WordCount = len(words)

	for _, word := range words {
		word = strings.ToLower(word)
		fi.TermFrequency[word]++
		i.FileFreq[word]++
	}

	i.FileInfo[filePath] = *fi
}

// Indexer method to remove file from index
func (i *Indexer) RemoveFile(filePath string) {
	// remove file from index
	if fi, exists := i.FileInfo[filePath]; exists {
		for term := range fi.TermFrequency {
			i.FileFreq[term] -= fi.TermFrequency[term]
			if i.FileFreq[term] <= 0 {
				delete(i.FileFreq, term)
			}
		}
		delete(i.FileInfo, filePath)
	}
}

// Indexer method to update file in index
func (i *Indexer) UpdateFile(filePath string) {
	// update file in index
	i.RemoveFile(filePath)
	i.AddFile(filePath)
}

// Indexer method to search for files which returns list of files
func (i *Indexer) Search(query string) []TermLoc {
	// search for files
	terms := strings.Fields(strings.ToLower(query))
	scores := make(map[string]float64)

	for _, term := range terms {
		idf := computeIDF(term, i)
		for filePath, fi := range i.FileInfo {
			tf := computeTF(term, &fi)
			scores[filePath] += tf * idf
		}
	}

	// fmt.Println(scores)

	selectedScores := make(map[string]float64)
	for filePath, score := range scores {
		if score > 0 {
			selectedScores[filePath] = score
		}
	}

	fmt.Println(selectedScores)

	// Collect results
	var results []TermLoc
	for filePath := range selectedScores {
		lines := []Loc{}

		// open file and search for lines
		fiSize := i.FileInfo[filePath].Size
		file, err := os.Open(filePath)
		if err != nil {
			continue
		}
		defer file.Close()

		data := make([]byte, fiSize)
		_, err = file.Read(data)
		if err != nil {
			continue
		}

		content := string(data)
		for i, line := range strings.Split(content, "\n") {
			for _, term := range terms {
				if strings.Contains(strings.ToLower(line), term) {
					lines = append(lines, Loc{Row: i, Col: strings.Index(strings.ToLower(line), term)})
				}
			}
		}

		results = append(results, TermLoc{FilePath: filePath, Lines: lines})
	}
	return results
}

// FileInfo constructor
func NewFileInfo() *FileInfo {
	return &FileInfo{
		TermFrequency: make(map[string]int),
	}
}

// function to compute tf
func computeTF(term string, fileInfo *FileInfo) float64 {
	// compute term frequency
	return float64(fileInfo.TermFrequency[term]) / float64(fileInfo.WordCount)
}

// function to compute idf
func computeIDF(term string, indexer *Indexer) float64 {
	// compute inverse document frequency
	totalDocs := len(indexer.FileInfo)
	docFreq := indexer.FileFreq[term]
	if docFreq == 0 {
		return 0
	}
	return math.Log(float64(totalDocs) / float64(docFreq))
}

func main() {
	fmt.Println("Hello, World!")

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <directory> <query>")
		return
	}

	// directory to index
	testdir := os.Args[1]

	// search query
	query := os.Args[2]

	// included file extensions
	extensions := []string{".txt", ".md"}

	// create indexer
	indexer := NewIndexer()

	// index files in directory
	indexer.IndexFiles(testdir, extensions)

	// fmt.Println(indexer)

	// search for files
	termLocs := indexer.Search(query)
	fmt.Println(termLocs)
}
