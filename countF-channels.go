package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type CountsResult struct {
	LineCount        int
	WordsCount       int
	VowelsCount      int
	PunctuationCount int
}

func Counts(chunk []byte, results chan<- CountsResult) {
    lineCount := 0
    wordsCount := 0
    vowelsCount := 0
    punctuationCount := 0

    inWord := false

    for _, b := range chunk {
        switch {
        case b == '\n':
            lineCount++
        case b == ' ' || b == '\t':
            if inWord {
                wordsCount++
                inWord = false
            }
        default:
            inWord = true

            if isVowel(b) {
                vowelsCount++
            }

            if isPunctuation(b) {
                punctuationCount++
            }
        }
    }

    if inWord {
        wordsCount++
    }

    results <- CountsResult{LineCount: lineCount, WordsCount: wordsCount, VowelsCount: vowelsCount, PunctuationCount: punctuationCount}
}


func isVowel(b byte) bool {
	vowels := "AEIOUaeiou"
	return byteInSlice(b, []byte(vowels))
}

func isPunctuation(b byte) bool {
	punctuation := "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	return byteInSlice(b, []byte(punctuation))
}

func byteInSlice(b byte, slice []byte) bool {
	for _, value := range slice {
		if b == value {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], " <file_path> <num_chunks>")
		return
	}

	filePath := os.Args[1]                     // Get the file path from the command-line argument
	numChunks, err := strconv.Atoi(os.Args[2]) // Get the number of chunks as an integer
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	results := make(chan CountsResult, numChunks)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSize := stat.Size()
	chunkSize := fileSize / int64(numChunks)

	reader := bufio.NewReader(file)

	for i := 0; i < numChunks; i++ {
		chunk := make([]byte, chunkSize)
		_, err := reader.Read(chunk)
		if err != nil {
			log.Fatal(err)
		}

		go Counts(chunk, results)
	}

	totalCounts := CountsResult{}

	for i := 0; i < numChunks; i++ {
		result := <-results
		totalCounts.LineCount += result.LineCount
		totalCounts.WordsCount += result.WordsCount
		totalCounts.VowelsCount += result.VowelsCount
		totalCounts.PunctuationCount += result.PunctuationCount
	}

	fmt.Println("Number of lines:", totalCounts.LineCount)
	fmt.Println("Number of words:", totalCounts.WordsCount)
	fmt.Println("Number of vowels:", totalCounts.VowelsCount)
	fmt.Println("Number of punctuation:", totalCounts.PunctuationCount)
	fmt.Println("Run Time:", time.Since(start))
}
