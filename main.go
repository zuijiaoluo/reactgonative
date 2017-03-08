package main

import (
	"fmt"
	"strings"

	"github.com/steve-winter/reactgonative/filebuilder"
	"github.com/steve-winter/reactgonative/goparser"
	"github.com/steve-winter/reactgonative/types"
)

var defaultAndroidRoot = "app/src/main/java/"
var defaultPackageRoot = "com.reactgohybrid"
var defaultGoPackage = "/golang.org/x/mobile/example/bind/hello"

func main() {
	fmt.Printf("Processing package %s\n", defaultGoPackage)
	tList, err := goparser.Parsing(defaultGoPackage)
	if err != nil {
		fmt.Printf("Unable to parse file - %s\n", err.Error())
	}
	for _, t := range tList {
		if t.IsValid() {
			fmt.Printf("\tPackagename created: %s\n", t.PackageName)
			typeString := module(t)
			err := packageBuild(typeString, t.PackageName)
			if err != nil {
				fmt.Printf("Unable to build package - %s\n", err.Error())
			}
		}
	}
}

func module(t types.GoType) string {
	m := filebuilder.NewModuleBuilder(defaultAndroidRoot,
		defaultPackageRoot)
	typeString, err := m.BuildModule(&t)
	if err != nil {
		fmt.Printf("Unable to build module - %s\n", err.Error())
		return ""
	}
	err = m.Close()
	if err != nil {
		fmt.Printf("Unable to build module - %s\n", err.Error())
		return ""
	}
	return typeString
}

func packageBuild(typeString string, packageName string) error {
	m := filebuilder.NewPackageBuilder(defaultAndroidRoot,
		defaultPackageRoot)

	err := m.BuildPackage(typeString, packageName)
	if err != nil {
		return err
	}
	err = m.Close()
	if err != nil {
		return err
	}
	return nil
}

func goToJavaType(javaType string) string {
	x := strings.ToLower(javaType)
	switch x {
	case "string":
		return "string"
	case "int":
		return "long"
	}
	return ""
}
