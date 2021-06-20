package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

//go:generate sqlgen -src ..\test\datatype.go:Customer

const (
	targetFile = `C:\Users\curtiscao\Documents\Workspace\sqlgen\tests\datatype.go`
)

func main() {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, targetFile, nil, parser.ParseComments)
	switch err.(type) {
	case nil:
	case scanner.ErrorList:
		fmt.Printf("There are syntax errors in your source code: %v\n", err)
	default:
		fmt.Printf("Failed to read your source code: %v\n", err)
		os.Exit(-1)
	}

	pre := func(cursor *astutil.Cursor) bool {
		if _, ok := cursor.Node().(ast.Decl); !ok {
			return false
		}
		if typ, ok := cursor.Node().(*ast.StructType); ok && cursor.Name() == "Node" {
			fmt.Println(typ)
		}
		return true
	}
	post := func(cursor *astutil.Cursor) bool {
		return true
	}
	astutil.Apply(f, pre, post)
}
