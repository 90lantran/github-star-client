package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/akamensky/argparse"
)

func main() {

	// Create new parser object
	parser := argparse.NewParser("main", "Sends a request and Receives a response from github-star server")

	// // give it verbose flag
	// verbose := parser.Flag(
	// 	"", "verbose", &argparse.Options{
	// 		Help: "Verbose mode",
	// 	},
	// )

	s := parser.StringList("r", "request", &argparse.Options{Required: true, Help: "List of organization/repossitory to send to server"})
	// Create string flag
	//s := parser.String("r", "request", &argparse.Options{Required: true, Help: "Json Payload to send to server"})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}
	// Finally print the collected string
	//fmt.Println(*s)
	a, _ := json.Marshal(*s)

	jsonStr := []byte(`{"input":` + fmt.Sprintf("%s", a) + `}`)

	fmt.Println(fmt.Sprintf("%s", jsonStr))

	req, err := http.NewRequest("POST", "http://localhost:8080/get-stars", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		//Timeout: 10000,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	log.Println(string(body))
	fmt.Printf("Response: %s\n", string(body))

}
