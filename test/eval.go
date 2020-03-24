package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/yhyddr/golangbot/eval"
)

var code = `
package main

import (
	"fmt"
)

func main() {
	%s
}
`

func main() {
	c := strings.NewReplacer(`“`, `"`, `”`, `"`).Replace(`fmt.Println(“Hello, World”)`)
	fmt.Println(c)
	// res, err := eval.GoCode(fmt.Sprintf(code, `fmt.Println("Hello, World")`))
	res, err := eval.GoCode(fmt.Sprintf(code, c))
	if err != nil {
		log.Println(err)
	}
	fmt.Println("error: ", res.Errors)
	fmt.Println("event: ", res.Events)
	fmt.Println("others: ", res.IsTest, res.Status, res.TestsFailed)
}
