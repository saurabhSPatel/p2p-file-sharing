package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadFile reads a file and returns its contents as a byte slice.
func ReadFile(path string) ([]byte, error) {
	cleanPath := filepath.Clean(path)
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %w", cleanPath, err)
	}
	return data, nil
}

// WriteFile saves the provided data to a specified path.
func WriteFile(path string, data []byte) error {
	cleanPath := filepath.Clean(path)
	if err := os.WriteFile(cleanPath, data, 0644); err != nil {
		return fmt.Errorf("error writing to file '%s': %w", cleanPath, err)
	}
	return nil
}

// ListFiles retrieves a list of all files in the specified directory.
func ListFiles(dir string) ([]string, error) {
	cleanDir := filepath.Clean(dir)
	var files []string
	err := filepath.Walk(cleanDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing file '%s': %w", path, err)
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory '%s': %w", cleanDir, err)
	}

	return files, nil
}
