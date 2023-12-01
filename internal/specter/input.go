package specter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func GetAbsolutePath(file string) (string, error) {
	// Get the absolute filepath
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	// Check if the file exists
	_, err = os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("%s does not exist", path)
	}

	return path, nil
}

func ReadFromFile(filename string) ([]byte, error) {
	content, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func PrintFileContent(content []byte) {
	fmt.Println(string(content[:]))
}

func ReadFromStdin() []byte {
	// Receive the input via stdin
	scanner := bufio.NewScanner(os.Stdin)
	var input []byte
	for scanner.Scan() {
		input = scanner.Bytes()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading stdin: ", err)
	}

	return input
}
