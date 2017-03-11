package goparser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/steve-winter/reactgonative/types"
)

const (
	gopath = "GOPATH"
)

//Parsing locates the Go package at pkgIdentifier and generates an array of GoType representing each package, or an error
func Parsing(pkgIdentifier string) ([]types.GoType, error) {
	pkgs, e := parsePackage(pkgIdentifier)
	if e != nil {
		return []types.GoType{}, e
	}
	typeList := make([]types.GoType, 0)
	for _, pkg := range pkgs {
		pkgName := pkg.Name
		for name, f := range pkg.Files {
			if !strings.HasSuffix(name, "_test.go") {
				typeList = append(typeList, parseFile(f, pkgName))
			}
		}
	}
	return typeList, nil
}

func parsePackage(pkgIdentifier string) (pkgs map[string]*ast.Package, first error) {
	folder := buildPackageFolder(pkgIdentifier)
	fset := token.NewFileSet()
	pkgs, e := parser.ParseDir(fset, folder, nil, 0)
	return pkgs, e
}

func buildPackageFolder(pkgIdentifier string) string {
	if len(os.Getenv(gopath)) == 0 {
		panic("GOPATH is not set")
	}
	folder := filepath.Join(os.Getenv(gopath), "src")
	folder = filepath.Join(folder, pkgIdentifier)
	return folder
}

func parseFile(f *ast.File, pkgName string) types.GoType {
	m := types.GoType{}
	m.PackageName = pkgName
	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			//Function declared
			parseFunc(x, &m)
			// case *ast.Package:
			// case *ast.FieldList:
			// case *ast.BasicLit:
			// case *ast.Ident:
			// case *ast.ReturnStmt:
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
			switch x := parameterList.Type.(type) {
			case *ast.Ident:
				paramsLength := len(m.Functions[len(m.Functions)-1].Params)
				m.Functions[len(m.Functions)-1].Params[paramsLength-1].T = x.Name
				for _, parameterName := range parameterList.Names {
					m.Functions[len(m.Functions)-1].Params[paramsLength-1].Name = parameterName.Name
				}
			}
		}
	}
}

func parseReturn(x *ast.FuncDecl, m *types.GoType) {
	m.Returns = append(m.Returns, types.GoParams{})
	if x.Type.Results != nil {
		for _, parameterList := range x.Type.Results.List {
			switch v := parameterList.Type.(type) {
			case *ast.Ident:
				// specialType := parameterList.Type.(*ast.Ident)
				m.Returns[len(m.Returns)-1].T = v.Name
				for _, parameterName := range parameterList.Names {
					m.Returns[len(m.Returns)-1].Name = parameterName.Name
				}
			case *ast.ArrayType:
				// fmt.Printf("ARRAYTYPE %s uuu %v\n", x.Name.String(), v.Elt.(*ast.SelectorExpr).Sel.Obj)

			}
		}
	}
}
