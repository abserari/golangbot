package main

// #cgo CFLAGS: -I${SRCDIR}
// #cgo LDFLAGS: -L${SRCDIR} -lhello
// #include "libhello.h"
import "C"
import "fmt"

func main() {
	fmt.Println("ing....")

	fmt.Println(C.Hello(10))
}
