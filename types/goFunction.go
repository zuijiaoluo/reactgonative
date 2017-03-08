package types

//GoFunction represents a Go functions name and an array of parameters, if any.
type GoFunction struct {
	Name   string
	Params []GoParams
}
