package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type VirtualScanner struct {
	DATE  string `xml:"DATETIME"`
	APPLS []struct {
		ID     string `xml:"ID"`
		NAME   string `xml:"NAME"`
		STATUS string `xml:"STATUS"`
	} `xml:"RESPONSE>APPLIANCE"`
}

// Function to fetch data from Qualys API
func fetchData(apiURL, username, password string) (*VirtualScanner, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var scanner VirtualScanner
	err = xml.Unmarshal(body, &scanner)
	if err != nil {
		return nil, err
	}

	return &scanner, nil
}

// Function to query the scanner by name
func queryScanner(scannerData *VirtualScanner, scannerName string) (string, string, error) {
	for _, appl := range scannerData.APPLS {
		if strings.EqualFold(appl.NAME, scannerName) {
			return appl.ID, appl.STATUS, nil
		}
	}
	return "", "", fmt.Errorf("scanner with name %s not found", scannerName)
}

func main() {
	apiURL := "https://qualysapi.qualys.com/api/2.0/fo/appliance/"
	username := "your_api_username"
	password := "your_api_password"

	// Fetch data from Qualys API
	scannerData, err := fetchData(apiURL, username, password)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}

	// Query the scanner by name
	scannerName := "Scanner Name"
	id, status, err := queryScanner(scannerData, scannerName)
	if err != nil {
		fmt.Println("Error querying scanner:", err)
	} else {
		fmt.Printf("Scanner ID: %s, Status: %s\n", id, status)
	}
}
