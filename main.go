package main

import "fmt"

func main() {
	for i := range 5 {
		fmt.Println(i / 13)
		fmt.Printf("page_%d\n", (i/13)+1)
	}
}
