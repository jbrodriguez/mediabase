package helper

import (
	"bufio"
	"os"
)

type StringMap map[string]string

func GetOrDefaultString(ask string, def string) string {
	if ask != "" {
		return ask
	}
	return def
}

func GetOrDefaultInt(ask int, def int) int {
	if ask != 0 {
		return ask
	}
	return def
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
