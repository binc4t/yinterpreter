package main

import (
	"fmt"
	"github.com/binc4t/yinterpreter/ast"
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
		parser := ast.NewParser(idy)
		prog := parser.ParseProgram()
		fmt.Println(prog.String())
		if len(parser.Errors()) != 0 {
			fmt.Println("err: ", parser.Errors())
		}

		//for t := idy.NextToken(); t.Type != identify.EOF; t = idy.NextToken() {
		//	fmt.Printf("%+v\n", t)
		//}
	}
}
