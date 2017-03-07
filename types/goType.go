package types

type GoType struct {
	PackageName string
	Functions   []GoFunction
	Returns     []GoParams
}

func (g *GoType) IsValid() bool {
	if g.PackageName != "" && len(g.Functions) > 0 {
		return true
	}
	return false
}
