package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func CIDRtoSubNet(cidr string) (string, error) {
	// Parse the CIDR notation
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	// Convert the subnet mask to dotted decimal format
	subnetMask := net.IP(ipnet.Mask).String()

	// Combine IP address and subnet mask
	return fmt.Sprintf("%s|%s", ip.String(), subnetMask), nil
}

func main() {
	// Example CIDR notation
	//cidr := "192.168.1.0/27"
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
	scanner := bufio.NewScanner(file)

	//slice to hold each line
	var lines []string
	// var FULL_
	//Read each line
	N := 0
	for scanner.Scan() {
		//fget a line from the file
		lines = append(lines, scanner.Text())
		// break the line into parts
		lineparts := strings.Split(lines[N], "|")
		if len(lineparts) != 3 {
			fmt.Println("Error: the line does not contain 3 parts")
			return
		}
		VLAN := lineparts[0]
		CIDR := lineparts[1]
		DESCR := lineparts[2]
		//fmt.Println(VLAN, ",", CIDR, ",", DESCR)

		// Convert to subnet mask notation
		subnet, err := CIDRtoSubNet(CIDR)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		//combine all the fields
		fmt.Print(VLAN, "|", subnet, "|", DESCR, "\n")

		N++
	}
	// Check for errors reading the line
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading from file: %s", err)
	}

}
