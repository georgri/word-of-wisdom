package quotes

import (
	"bufio"
	"os"
	"strings"
)

const (
	localFileWithQuotes = "data/quotes.txt"
)

// GetHardcodedQuotes returns a predefined hardcoded quote list
func GetHardcodedQuotes() []string {
	quotes, err := GetQuotesFromFile(localFileWithQuotes)
	if err != nil {
		panic(err)
	}
	return quotes
}

// GetQuotesFromFile reads a quote list from a given plain text file
func GetQuotesFromFile(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
