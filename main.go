package main

import (
	"fmt"
	"github.com/binc4t/yinterpreter/identify"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need file path")
		return
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	idy := identify.NewIdentifier(f)

	for idy.FillIn() {
		for t := idy.NextToken(); t.Type != identify.EOF; t = idy.NextToken() {
			fmt.Printf("%+v\n", t)
		}
	}
}
