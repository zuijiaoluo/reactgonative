package main

import (
	"strings"

	"github.com/steve-winter/loggers"
	"github.com/steve-winter/reactgonative/filebuilder"
	"github.com/steve-winter/reactgonative/goparser"
	"github.com/steve-winter/reactgonative/types"
)

func main() {
	loggers.Info("Starting")
	tList := goparser.Parsing()
	for _, t := range tList {
		typeString := module(t)
		packageBuild(typeString, t.PackageName)
	}
}

func module(t types.GoType) string {
	m := filebuilder.NewModuleBuilder("/Users/SteveWinter/Development/golang/src/github.com/steve-winter/reactgonative/generated/test1",
		"com.reactgohybrid")
	err := m.Create()
	if err != nil {
		loggers.Errorf("Error %v", err)
	}
	typeString, err := m.BuildModule(&t)
	if err != nil {
		loggers.Errorf("Error %v", err)
	}
	err = m.Close()
	if err != nil {
		loggers.Errorf("Error %v", err)
	}
	return typeString
}

func packageBuild(typeString string, packageName string) string {
	m := filebuilder.NewPackageBuilder("/Users/SteveWinter/Development/golang/src/github.com/steve-winter/reactgonative/generated/testpackage1",
		"com.reactgohybrid")
	err := m.Create()
	if err != nil {
		loggers.Errorf("Error %v", err)
	}
	err = m.BuildPackage(typeString, packageName)
	if err != nil {
		loggers.Errorf("Error %v", err)
	}
	err = m.Close()
	if err != nil {
		loggers.Errorf("Error %v", err)
	}
	return typeString
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
