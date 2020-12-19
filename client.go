package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/akamensky/argparse"
)

type Request struct {
	Input []string `json:"input"`
}

type Response struct {
	TotalStars   int64         `json:"totalStars,omitempty"`
	InvalidRepos []string      `json:"invalidRepos,omitempty"`
	ValidRepos   []MapNameStar `json:"validRepos,omitempty"`
	Error        string        `json:"error,omitempty"`
	Status       string        `json:"status" validate:"required"`
}

type MapNameStar struct {
	Name string `json:"name"`
	Star int64  `json:"star(s)"`
}

func validateListInput(list []string) error {
	var validInput = regexp.MustCompile(`^[a-zA-Z0-9\_\-\.]+\/[a-zA-Z0-9\_\-\.]+$`)
	for _, l := range list {
		elements := strings.Split(l, ",")
		for _, e := range elements {
			if !validInput.MatchString(strings.TrimSpace(e)) {
				return fmt.Errorf("-r input list %s is not valid. Valid format is list of organization/repository", e)
			}
		}
	}
	return nil
}

func validateHostAndPort(hosts []string) error {
	var validHostAndPort = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):[0-9]+$`)
	for _, host := range hosts {
		if !validHostAndPort.MatchString(host) {
			return fmt.Errorf("-t Host and Port %s are not valid. Valid format ip:port", host)
		}
	}
	return nil
}

var argList *[]string
var argHost *string
var parser *argparse.Parser

func init() {
	//create create parser
	parser = argparse.NewParser("client", "sends a request and receives a response from github-star server")

	argList = parser.StringList("r", "request",
		&argparse.Options{
			Required: true,
			Help:     "List of organization/repossitory to send to server.",
			Validate: validateListInput})

	argHost = parser.String("t", "host",
		&argparse.Options{
			Required: false,
			Help:     "IP address and port of the server.",
			Validate: validateHostAndPort,
			Default:  "localhost:8080"})
}

func main() {
	// parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// create list of orginization/repository
	completeList := []string{}
	for _, arg := range *argList {
		elements := strings.Split(arg, ",")
		completeList = append(completeList, elements...)
	}

	request := Request{
		Input: completeList,
	}
	fmt.Printf("input list: %+v\n", request)

	jsonStr, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", "http://"+*argHost+"/get-stars", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatalf("cannot create a POST request %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// create default http client
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client cannot send request %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body, &response)
	// pretty print json response
	responseJSON, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("Response: %+v\n", string(responseJSON))
}
