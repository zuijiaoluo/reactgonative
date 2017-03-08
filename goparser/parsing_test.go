package goparser

import (
	"go/ast"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/steve-winter/reactgonative/types"
)

func TestParseFuncName(t *testing.T) {
	Convey("Given empty go type", t, func() {
		m := &types.GoType{}

		Convey("When a function is added", func() {
			var x ast.FuncDecl
			x.Name = &ast.Ident{}
			x.Name.Name = "functionName1"
			parseFuncName(&x, m)

			Convey("Then the length should be 1", func() {
				So(len(m.Functions), ShouldEqual, 1)
			})
		})
	})
}
