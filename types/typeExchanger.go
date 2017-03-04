package types

import "strings"

func GoToJava(goIn string) string {
	x := strings.ToLower(goIn)
	switch x {
	case "string":
		return "String"
	}
	return goIn
}

func JavaToGo(javaIn string) string {

	return javaIn
}
