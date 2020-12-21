package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationInputList(t *testing.T) {
	list := []string{"me/"}
	err := ValidateListInput(list)
	assert.NotNil(t, err)

	list = []string{"me/me"}
	err = ValidateListInput(list)
	assert.Nil(t, err)
}

func TestValidationHostAndPort(t *testing.T) {
	host := []string{"http://127.0.0.1"}
	err := ValidateHostAndPort(host)
	assert.NotNil(t, err)

	host = []string{"http://127.0.0.1:8080"}
	err = ValidateHostAndPort(host)
	assert.Nil(t, err)
}

func TestFunctionalityOfClient(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Contains(t, req.URL.String(), "/get-stars")
		mapNameStar := MapNameStar{
			Name: "golang/go",
			Star: 5,
		}
		payload := Payload{
			TotalStars:   5,
			InvalidRepos: []string{"golang/o"},
			ValidRepos:   []MapNameStar{mapNameStar},
		}
		userResponse := UserResponse{
			Pl:     &payload,
			Error:  "no error",
			Status: "good",
		}
		jsonData, err := json.Marshal(userResponse)
		if err != nil {
			fmt.Printf("cannot marshal userResponse %v\n", err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(jsonData)
	}))
	// Close the server when test finishes
	defer server.Close()

	fmt.Println("server.URL: " + server.URL)
	fmt.Printf("server.Client: %+v\n", server.Client())

	argList := []string{"golang/go", "golang/o"}
	argHost := server.URL

	if err := ValidateListInput(argList); err != nil {
		fmt.Printf("invalid list input %v\n", err)
		return
	}

	if err := ValidateHostAndPort([]string{argHost}); err != nil {
		fmt.Printf("invalid host and port %v\n", err)
		return
	}

	req, err := CreatePostRequest(&argList, &argHost)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("req: %+v\n", req)

	// Use Client & URL from our local test server
	customizedClient := CustomizedClient{server.Client(), &server.URL}
	resp, err := customizedClient.SendPostRequest(req)

	err = ShowResponse(resp)
	assert.Nil(t, err)

}
