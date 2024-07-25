package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	//make sure filename is given on command line
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <filename>", os.Args[0])
	}
	// Get the filename from the cmd line
	filename := os.Args[1]
	
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()
	
	//screate a new scanner for the file
	scanner := bufio.NewScanner(file)
	
	//slice to hold each line
	var lines []string
	
	//Read each line
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		lines = append(lines, scanner.Text())
	}
	
	// Check for errors reading the line
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading from file: %s",err)
	}
	// join the lines with a comma
	result := strings.Join(lines, ",")
	
	fmt.Println(result)
}