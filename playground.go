package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Events struct {
	Message string
	Kind    string
	Delay   int
}

type GolangResult struct {
	Errors      string
	Events      []Events
	Status      int
	IsTest      bool
	TestsFailed int
}

func main() {
	var r http.Request
	r.ParseForm()
	r.Form.Add("version", "2")
	r.Form.Add("body", `package main

	import (
		"fmt"
	)
	
	func main() {
		fmt.Println("Hello, playground")
	}`)
	body := strings.NewReader(r.Form.Encode())
	resp, err := http.Post("https://play.golang.org/compile", "application/x-www-form-urlencoded", body)
	if err != nil {
		fmt.Println(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var res GolangResult
	err = json.Unmarshal(content, &res)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res.Events)
}
