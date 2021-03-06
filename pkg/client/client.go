package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Pl     *Payload `json:"payload,omitempty"`
	Error  []string `json:"error,omitempty"`
	Status string   `json:"status" validate:"required"`
}
type Payload struct {
	TotalStars   int64         `json:"totalStars,omitempty"`
	InvalidRepos []string      `json:"invalidRepos,omitempty"`
	ValidRepos   []MapNameStar `json:"validRepos,omitempty"`
}

type MapNameStar struct {
	Name string `json:"name"`
	Star int64  `json:"star(s)"`
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
		return nil, err
	}
	req, err = http.NewRequest("POST", *baseURL+"/get-stars", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
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

func ShowResponse(resp *http.Response) (*UserResponse, error) {
	var userResponse UserResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}

	// pretty print json response
	responseJSON, err := json.MarshalIndent(userResponse, "", "  ")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Response: %+v\n", string(responseJSON))
	return &userResponse, nil
}
