package goparser

import (
	"go/ast"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/steve-winter/reactgonative/types"
)

func TestParseFuncName(t *testing.T) {
	Convey("Given empty go type", t, func() {
		m := &types.GoType{}
		Convey("When a function is added", func() {
			var x ast.FuncDecl
			x.Name = &ast.Ident{}
			x.Name.Name = "functionName1"
			parseFuncName(&x, m)
			Convey("Then the length should be 1", func() {
				So(len(m.Functions), ShouldEqual, 1)
			})
		})
		Convey("When a second function is added", func() {
			y := ast.FuncDecl{
				Name: &ast.Ident{
					Name: "functionName1",
				},
			}
			parseFuncName(&y, m)
			x := ast.FuncDecl{
				Name: &ast.Ident{
					Name: "functionName2",
				},
			}
			parseFuncName(&x, m)
			Convey("Then the length should be 2", func() {
				So(len(m.Functions), ShouldEqual, 2)
			})
		})
	})
}

func TestBuildPackageFolder(t *testing.T) {
	realGoPath := os.Getenv("GOPATH")
	Convey("Given gopath is set", t, func() {
		Convey("When a package name ok pkg is used", func() {
			folder := buildPackageFolder("pkg")
			Convey("Then the folder name should be valid", func() {
				So(folder, ShouldEqual, filepath.Join(realGoPath, "src/pkg"))
			})
		})
	})

	Convey("Given gopath is not set", t, func() {
		os.Setenv("GOPATH", "")
		Convey("When a package name pkg is used", func() {
			Convey("Then the folder name should panic", func() {
				So(func() { buildPackageFolder("pkg") }, ShouldPanicWith, "GOPATH is not set")
			})
		})
		os.Setenv("GOPATH", realGoPath)
	})
}

func TestParseFile(t *testing.T) {
	Convey("Given empty file", t, func() {
		file := &ast.File{}
		Convey("When only a package name pkg is used", func() {
			goType := parseFile(file, "pkg")
			Convey("Then the package name on the gotype is pkg", func() {
				So(goType.PackageName, ShouldEqual, "pkg")
			})
			Convey("Then the number of functions should be 0", func() {
				So(len(goType.Functions), ShouldEqual, 0)
			})
			Convey("Then the number of function returns should be 0", func() {
				So(len(goType.Returns), ShouldEqual, 0)
			})
		})
	})

	Convey("Given a file with functions", t, func() {
		file := &ast.File{
			Decls: []ast.Decl{
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "FunctionName1",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								&ast.Field{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "Param1",
										},
									},
									Type: &ast.Ident{
										Name: "string",
									},
								},
							},
						},
					},
				},
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "nonExported",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								&ast.Field{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "param1",
										},
									},
									Type: &ast.Ident{
										Name: "int",
									},
								},
							},
						},
					},
				},
				&ast.FuncDecl{
					Name: &ast.Ident{
						Name: "ExportedReturns",
					},
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								&ast.Field{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "param1",
										},
									},
									Type: &ast.Ident{
										Name: "int",
									},
								},
							},
						},
						Results: &ast.FieldList{
							List: []*ast.Field{
								&ast.Field{
									Names: []*ast.Ident{
										&ast.Ident{
											Name: "return1",
										},
									},
									Type: &ast.Ident{
										Name: "int",
									},
								},
							},
						},
					},
				},
			},
		}
		Convey("When an exported and unexported function is set", func() {
			goType := parseFile(file, "pkg")
			Convey("Then the number of functions equals 1", func() {
				So(len(goType.Functions), ShouldEqual, 2)
			})
			Convey("And has 1 parameter", func() {
				So(len(goType.Functions[0].Params), ShouldEqual, 1)
			})
			Convey("And has function name FunctionName1", func() {
				So(goType.Functions[0].Name, ShouldEqual, "FunctionName1")
			})
			Convey("And has function param name Param1", func() {
				So(goType.Functions[0].Params[0].Name, ShouldEqual, "Param1")
			})
			Convey("And one has a return type of int", func() {
				So(goType.Returns[1].T, ShouldEqual, "int")
				Convey("And name of return1", func() {
					So(goType.Returns[1].Name, ShouldEqual, "return1")
				})
			})
		})
	})
}

func TestParsePackage(t *testing.T) {
	Convey("Given a package directory of the goparser", t, func() {
		pkgDir := "github.com/steve-winter/reactgonative/goparser"
		Convey("When parse package is called", func() {
			pkgs, err := parsePackage(pkgDir)
			Convey("Then there are no errors", func() {
				So(err, ShouldBeNil)
			})
			Convey("And there is 1 package found", func() {
				So(len(pkgs), ShouldEqual, 1)
			})
			Convey("And there are 2 files found", func() {
				So(len(pkgs["goparser"].Files), ShouldEqual, 2)
			})
			Convey("And there are 10 declarations", func() {
				fileName := filepath.Join(os.Getenv("GOPATH"), "src")
				fileName = filepath.Join(fileName, pkgDir)
				fileName = filepath.Join(fileName, "parsing.go")
				So(len(pkgs["goparser"].Files[fileName].Decls), ShouldEqual, 10)
			})
		})
	})
}

func TestParsing(t *testing.T) {
	Convey("Given a package directory of the goparser", t, func() {
		pkgDir := "github.com/steve-winter/reactgonative/goparser"
		Convey("When parsing is called", func() {
			goTypes, err := Parsing(pkgDir)
			Convey("Then there are no errors", func() {
				So(err, ShouldBeNil)
			})
			Convey("And there is 1 package found", func() {
				So(len(goTypes), ShouldEqual, 1)
			})
			Convey("And 1 exported function", func() {
				So(len(goTypes[0].Functions), ShouldEqual, 1)
			})
			Convey("And function name of Parsing", func() {
				So(goTypes[0].Functions[0].Name, ShouldEqual, "Parsing")
			})
			Convey("And the function has 0 params", func() {
				So(len(goTypes[0].Functions[0].Params), ShouldEqual, 1)
			})
			Convey("And with a param name of pkgIdentifier with type of string", func() {
				So(goTypes[0].Functions[0].Params[0].Name, ShouldEqual, "pkgIdentifier")
				So(goTypes[0].Functions[0].Params[0].T, ShouldEqual, "string")
			})
			Convey("And return type is error (BUG)", func() {
				//BUG - Issue #9. Only supports single return type at present
				So(goTypes[0].Returns[0].T, ShouldEqual, "error")
			})
		})
	})
	Convey("Given a package doesnt exist", t, func() {
		Convey("When parsing is called", func() {
			pkgs, err := Parsing("384738h932h392h32")
			Convey("Then an error is returned", func() {
				So(err.Error(), ShouldEndWith, "no such file or directory")
			})
			Convey("And no packages returned", func() {
				So(len(pkgs), ShouldEqual, 0)
			})
		})
	})
}
