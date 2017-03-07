package main

import (
	"strings"

	"github.com/steve-winter/loggers"
	"github.com/steve-winter/reactgonative/filebuilder"
	"github.com/steve-winter/reactgonative/goparser"
	"github.com/steve-winter/reactgonative/types"
)

var defaultAndroidRoot = "app/src/main/java/"
var defaultPackageRoot = "com.reactgohybrid"

func main() {
	loggers.Info("Starting")
	tList := goparser.Parsing()
	for _, t := range tList {
		if t.IsValid() {
			loggers.Infof("Package name %s", t.PackageName)
			typeString := module(t)
			packageBuild(typeString, t.PackageName)
		}
	}
}

func module(t types.GoType) string {
	m := filebuilder.NewModuleBuilder(defaultAndroidRoot,
		defaultPackageRoot)
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

func packageBuild(typeString string, packageName string) (string, error) {
	m := filebuilder.NewPackageBuilder(defaultAndroidRoot,
		defaultPackageRoot)
	err := m.BuildPackage(typeString, packageName)
	if err != nil {
		loggers.Errorf("Error %v", err)
		return "", err
	}
	err = m.Close()
	if err != nil {
		loggers.Errorf("Error %v", err)
	}
	return typeString, nil
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
