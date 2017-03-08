package types

//GoType represents a single Go file
type GoType struct {
	PackageName string
	Functions   []GoFunction
	Returns     []GoParams
}

//IsValid identifies whether the GoType holds valid data
func (g *GoType) IsValid() bool {
	if g.PackageName != "" && len(g.Functions) > 0 {
		return true
	}
	return false
}
