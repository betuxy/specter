package specter

import (
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
