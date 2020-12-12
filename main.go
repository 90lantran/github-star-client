package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	var jsonStr = []byte(`{"input":["tinygo-org/tinygo-site",
	"tinygo-org/homebrew-tools"]}`)
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

	log.Println("Me: ")
	log.Println(string(body))
	fmt.Printf("Response: %s\n", string(body))

}
