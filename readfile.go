package main

import (
	"fmt"
)

func readFile(filename string) (string, error) {
	// Read the entire file content
	content, err := io.util.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// Convert the content to a string and return
	return string(content), nil
}

func main() {
	// Example usage
	filename := "example.txt"
	content, err := readFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("File content:")
	fmt.Println(content)
}
