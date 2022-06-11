package main

import "C"

func main() {
}

//export Hello
func Hello(in int) int {
	var ch chan struct{}
	for i := 0; i < 100; i++ {
		go func() {
			in += 1
		}()
	}
	ch <- struct{}{}

	return in
}
