package helpers

import "fmt"

func Trace(s string) string {
	fmt.Println("Entering:", s)
	return s
}

func Un(s string) {
	fmt.Println("Leaving:", s)
}
