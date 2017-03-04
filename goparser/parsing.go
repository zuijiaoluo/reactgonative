package goparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/steve-winter/loggers"
	"github.com/steve-winter/reactgonative/types"
)

func Parsing() []types.GoType {
	folder := filepath.Join(os.Getenv("GOPATH"), "src")
	folder = filepath.Join(folder, "/golang.org/x/mobile/example/bind/hello")
	fset := token.NewFileSet()
	pkgs, e := parser.ParseDir(fset, folder, nil, 0)
	if e != nil {
		loggers.Infof("Filename: %s", folder)
		loggers.Fatal(e)
		return []types.GoType{}
	}
	typeList := make([]types.GoType, 1)
	for _, pkg := range pkgs {
		pkgName := pkg.Name
		for _, f := range pkg.Files {
			typeList = append(typeList, parseFile(fset, f, pkgName))
		}
	}
	return typeList
}

func parseFile(fset *token.FileSet, f *ast.File, pkgName string) types.GoType {
	ast.Print(fset, f)
	m := types.GoType{}
	m.PackageName = pkgName
	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			//Function declared
			parseFunc(x, &m)
			fmt.Printf("Function: %v\n\n", m)

		case *ast.Package:
		case *ast.FieldList:
		case *ast.BasicLit:
		case *ast.Ident:
		case *ast.ReturnStmt:
		}
		return true
	})

	// for _, n := range f.Decls {
	// 	fmt.Printf("Sub1 %v\n", n.Pos())
	// }
	return m
}

func parseFunc(x *ast.FuncDecl, m *types.GoType) {
	if x.Name.IsExported() {
		parseFuncName(x, m)
		parseParams(x, m)
		parseReturn(x, m)
	} else {
		loggers.Infof("Function not exported \"%s\".\"%s\"\n", m.PackageName, x.Name.Name)
	}
}

func parseFuncName(x *ast.FuncDecl, m *types.GoType) {
	functionName := x.Name.String()
	m.Functions = append(m.Functions, types.GoFunction{
		Name: functionName,
	})
}

func parseParams(x *ast.FuncDecl, m *types.GoType) {
	if x.Type.Params != nil {
		for _, parameterList := range x.Type.Params.List {
			m.Functions[len(m.Functions)-1].Params = append(m.Functions[len(m.Functions)-1].Params,
				types.GoParams{})
			specialType := parameterList.Type.(*ast.Ident)
			paramsLength := len(m.Functions[len(m.Functions)-1].Params)
			m.Functions[len(m.Functions)-1].Params[paramsLength-1].T = specialType.Name
			for _, parameterName := range parameterList.Names {
				m.Functions[len(m.Functions)-1].Params[paramsLength-1].Name = parameterName.Name
			}
		}
	}
}

func parseReturn(x *ast.FuncDecl, m *types.GoType) {
	m.Returns = append(m.Returns, types.GoParams{})
	if x.Type.Results != nil {
		for _, parameterList := range x.Type.Results.List {
			specialType := parameterList.Type.(*ast.Ident)
			m.Returns[len(m.Returns)-1].T = specialType.Name
			for _, parameterName := range parameterList.Names {
				m.Returns[len(m.Returns)-1].Name = parameterName.Name
			}
		}
	}
}
