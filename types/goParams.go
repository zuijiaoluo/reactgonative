package types

//GoParams represents an individual Go functions return type, or parameters.
//Name can be blank
type GoParams struct {
	Name string
	T    string
}
