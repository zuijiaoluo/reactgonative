package types

import "strings"

//GoToJava converts the goIn Go type, to the Java representation.
//BUG - Unfinished
func GoToJava(goIn string) string {
	x := strings.ToLower(goIn)
	switch x {
	case "string":
		return "String"
	}
	return goIn
}

//JavaToGo converts the javaIn Java type, to the Go representation
//BUG - Unfinished
func JavaToGo(javaIn string) string {

	return javaIn
}
