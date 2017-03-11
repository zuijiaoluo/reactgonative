package config

import "flag"

var flagvar int

func init() {
	flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
}

func ParseFlags() {
	if !flag.Parsed() {
		flag.Parse()
	}
}
