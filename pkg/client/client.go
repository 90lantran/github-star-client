package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type CustomizedClient struct {
	Client  *http.Client
	BaseURL *string
}

type UserRequest struct {
	Input []string `json:"input"`
}

type UserResponse struct {
	TotalStars   int64         `json:"totalStars,omitempty"`
	InvalidRepos []string      `json:"invalidRepos,omitempty"`
	ValidRepos   []MapNameStar `json:"validRepos,omitempty"`
	Error        string        `json:"error,omitempty"`
	Status       string        `json:"status"`
}

type MapNameStar struct {
	Name string `json:"name,omitempty"`
	Star int64  `json:"star(s),omitempty"`
}

func ValidateListInput(list []string) error {
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

func ValidateHostAndPort(hosts []string) error {
	var validHostAndPort = regexp.MustCompile(`^http:\/\/(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):[0-9]+$`)
	for _, host := range hosts {
		if !validHostAndPort.MatchString(host) {
			return fmt.Errorf("-t Host and Port %s are not valid. Valid format http://ip:port. If you dont specify it, http://localhost:8080 will be used", host)
		}
	}
	return nil
}

func CreatePostRequest(argList *[]string, baseURL *string) (req *http.Request, err error) {
	// create list of orginization/repository
	completeList := []string{}
	for _, arg := range *argList {
		elements := strings.Split(arg, ",")
		completeList = append(completeList, elements...)
	}

	userRequest := UserRequest{
		Input: completeList,
	}
	fmt.Printf("input list: %+v\n", userRequest)

	jsonStr, err := json.Marshal(userRequest)
	if err != nil {
		log.Printf("cannot marshall request %v\n", err)
		return
	}
	req, err = http.NewRequest("POST", *baseURL+"/get-stars", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Printf("cannot create a POST request %v\n", err)
		return

	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *CustomizedClient) SendPostRequest(req *http.Request) (resp *http.Response, err error) {
	resp, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ShowResponse(resp *http.Response) (err error) {
	var userResponse UserResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return err
	}
	// pretty print json response
	responseJSON, err := json.MarshalIndent(userResponse, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("Response: %+v\n", string(responseJSON))
	return nil
}
