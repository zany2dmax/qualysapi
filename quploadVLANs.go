package main

import (
	"encoding/base64"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type VirtualScanners struct {
	//super important to get the XML response fields properly nested or you wont get data in
	// these data structures
	XMLName  xml.Name         `xml:"APPLIANCE_LIST_OUTPUT"`
	Scanners []VirtualScanner `xml:"RESPONSE>APPLIANCE_LIST>APPLIANCE"`
}

type VirtualScanner struct {
	ID     string `xml:"ID"`
	Name   string `xml:"NAME"`
	Status string `xml:"STATUS"`
}

func Get_Credential_Hash(User string, Password string) string {

	return base64.StdEncoding.EncodeToString([]byte(User + ":" + Password))
}

func Get_Command_Line_Args() (string, string, string, string, string) {
	/* Get cmd line paramters */
	UserPtr := flag.String("username", "cmerc_el2", "Qualys Account User Name")
	PasswordPtr := flag.String("password", "C0m3r1caSucks!", "Qualys Account password")
	APIURLPtr := flag.String("APIURL", "https://qualysapi.qualys.com/", "Qualys API endpoint")
	scannerName := flag.String("scannerName", "QualysEGR1", "ScannerName")
	VLANFile := flag.String("filename", "", "VLAN filename")
	flag.Parse()
	return *UserPtr, *PasswordPtr, *APIURLPtr, *scannerName, *VLANFile
}

func Usage() {
	fmt.Println("usage:  quploadVLANs [-user -password -APIUrl -scannername -filename]")
	fmt.Println("     -scannername is required, must match a Qualys Scanner name")
	fmt.Println("     -filename is required,  must be a text file with one VLAN entry")
	fmt.Println("				per line in the form VLAN|NETWORK|SUBNET MASK|COMMENT")
}

// Function to fetch data from Qualys API
func listScannerData(apiURL, EncodedCred string) (*VirtualScanners, error) {
	// Setup the call to the Appliance API and return the array of scanners from the XML
	resource := "api/2.0/fo/appliance/"
	data := url.Values{}
	data.Set("action", "list")
	data.Add("output_mode", "brief")
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)
	fmt.Println("Calling Applance API ", urlStr)
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	req.Header.Add("X-requested-With", "GOLANG")
	req.Header.Add("authorization", "Basic "+EncodedCred)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	//fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected HTTP status: %s", resp.Status)
	} else {
		fmt.Println("Appliance API successful return....")
		//fmt.Println(resp.StatusCode)
	}
	//Setup the XML for querying the scannerName
	var vscanners VirtualScanners
	err = xml.Unmarshal(body, &vscanners)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))

	return &vscanners, nil
}

// Function to query the scanner by name
func queryScanner(scanners *VirtualScanners, scannerName string) (string, string, error) {
	for _, scanner := range scanners.Scanners {
		//fmt.Println("%s\n", scanner.Name)
		if strings.EqualFold(scanner.Name, scannerName) {
			return scanner.ID, scanner.Status, nil
		}
	}
	return "", "", fmt.Errorf("scanner with name %s not found", scannerName)
}

func readVLANfile(filename string) (string, error) {
	// Read the entire file content
	// Noet this file must be in the form:
	// VLAN|SUBNET|NETMASK|NAME,VLAN|SUBNET|NETMASK|NAME,...
	// where ...can be any number of VLANS separated by a comma
	// the set_vlans value for the API is an all or nothing parameter
	// in the setupVLANs function
	//
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	//fmt.Println(content, "\n", string(content))
	// Convert the content to a string and return
	return string(content), nil
}

func setupVLANs(apiURL string, ID string, STATUS string, VLANS string, EncodedCred string) error {
	// as noted in readVLANfile() format must be exact, and concatenated together with a ','
	resource := "api/2.0/fo/appliance/"
	data := url.Values{}
	data.Set("action", "update")
	data.Add("id", ID)
	fmt.Println(VLANS)
	data.Add("set_vlans", VLANS)
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)
	fmt.Println("Calling Applance API ", urlStr)
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
	}
	req.Header.Add("X-requested-With", "GOLANG")
	req.Header.Add("authorization", "Basic "+EncodedCred)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println("Response: ", resp.StatusCode, ":", string(body))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected HTTP status: %s", resp.Status)
	} else {
		fmt.Println("Appliance API update successful return....%s", body)
		//fmt.Println(resp.StatusCode)
	}
	return nil
}

func main() {
	username, password, apiURL, scannerName, filename := Get_Command_Line_Args()
	EncodedCred := Get_Credential_Hash(username, password)

	// Fetch data from Qualys Scanner API list
	vscanners, err := listScannerData(apiURL, EncodedCred)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	// Query the scanner by name
	id, status, err := queryScanner(vscanners, scannerName)
	if err != nil {
		fmt.Println("Error querying scanner:", err)
		os.Exit(0)
	} else {
		fmt.Printf("Scanner ID: %s, Status: %s\n", id, status)
	}
	//Now we have the ScannerID to setup the VLANs
	VLANS, err := readVLANfile(filename)
	if err != nil {
		fmt.Println("Error fetching file %s:", filename, err)
		return
	}
	fmt.Println(VLANS)
	//os.Exit(0)
	setupVLANs(apiURL, id, status, VLANS, EncodedCred)

}
