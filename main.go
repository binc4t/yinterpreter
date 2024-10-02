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
		for _, s := range prog.Statements {
			if l, ok := s.(*ast.LetStatement); ok {
				fmt.Println(s.TokenRaw(), l.Left.TokenRaw(), l.Right.TokenRaw())
			}
		}
		fmt.Println("err: ", parser.Errors())

		//for t := idy.NextToken(); t.Type != identify.EOF; t = idy.NextToken() {
		//	fmt.Printf("%+v\n", t)
		//}
	}
}
