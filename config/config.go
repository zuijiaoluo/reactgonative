package config

import "flag"

var flagvar int

func init() {
	flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
}

//ParseFlags processes flags for application. Should be called in main
func ParseFlags() {
	if !flag.Parsed() {
		flag.Parse()
	}
}
