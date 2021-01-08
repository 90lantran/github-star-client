package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/akamensky/argparse"

	. "github.com/90lantran/github-star-client/pkg/client"
)

func main() {
	//create create parser
	parser := argparse.NewParser("client", "sends a request and receives a response from github-star server")

	argList := parser.StringList("r", "request",
		&argparse.Options{
			Required: true,
			Help:     "List of organization/repossitory to send to server.",
			Validate: ValidateListInput})

	argHost := parser.String("t", "host",
		&argparse.Options{
			Required: false,
			Help:     "IP address and port of the server",
			Validate: ValidateHostAndPort,
			Default:  "http://localhost:8080"})

	// parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	customizedClient := CustomizedClient{
		Client:  &http.Client{},
		BaseURL: argHost,
	}

	req, err := CreatePostRequest(argList, argHost)
	if err != nil {
		fmt.Printf("cannot create POST request %v\n", err)
		return
	}
	resp, err := customizedClient.SendPostRequest(req)
	if err != nil {
		fmt.Printf("server is not up: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if _, err = ShowResponse(resp); err != nil {
		fmt.Printf("cannot show response from server %v\n", err)
		return
	}
}
