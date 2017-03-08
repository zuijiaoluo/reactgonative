package main

import "testing"

func TestGoToJavaType(t *testing.T) {
	if goToJavaType("string") != "string" {
		t.Fail()
	}
}
