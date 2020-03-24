package eval

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const helloworld = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}`

// Events for JSON
type Events struct {
	Message string
	Kind    string
	Delay   int
}

// GolangResult for JSON
type GolangResult struct {
	Errors      string
	Events      []Events
	Status      int
	IsTest      bool
	TestsFailed int
}

// GoCode run code and response
func GoCode(code string) (*GolangResult, error) {
	var r http.Request
	r.ParseForm()
	r.Form.Add("version", "2")
	if code == "" {
		code = helloworld
	}
	r.Form.Add("body", code)
	body := strings.NewReader(r.Form.Encode())
	resp, err := http.Post("https://play.golang.org/compile", "application/x-www-form-urlencoded", body)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}
	var res GolangResult
	err = json.Unmarshal(content, &res)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	return &res, nil
}
