package filebuilder

import (
	"strings"

	"github.com/steve-winter/loggers"
	"github.com/steve-winter/reactgonative/types"
)

var context = "reactContext"

type ModuleBuilder struct {
	javaFile *JavaFile
}

func NewModuleBuilder(name string, root string) ModuleBuilder {
	return ModuleBuilder{
		javaFile: NewJavaFile(name, root),
	}
}

func (mb *ModuleBuilder) BuildModule(g *types.GoType) (string, error) {
	err := mb.javaFile.WritePackageLine(g.PackageName)
	if err != nil {
		return "", err
	}
	err = mb.javaFile.WriteBlank(1)
	if err != nil {
		return "", err
	}
	err = mb.BuildImports(g)
	if err != nil {
		return "", err
	}

	err = mb.javaFile.WriteClassHeader(mb.className(g.PackageName),
		"ReactContextBaseJavaModule", "")
	if err != nil {
		return "", err
	}
	err = mb.javaFile.WriteBlank(1)
	if err != nil {
		return "", err
	}
	err = mb.BuildConstructor(g)
	if err != nil {
		return "", err
	}

	err = mb.BuildGetName(g)
	if err != nil {
		return "", err
	}
	err = mb.javaFile.WriteBlank(1)
	if err != nil {
		return "", err
	}

	err = mb.BuildReactMethods(&g.Functions, &g.Returns, g.PackageName)
	if err != nil {
		return "", err
	}
	err = mb.javaFile.WriteCloseTag()
	if err != nil {
		return "", err
	}
	return mb.className(g.PackageName), nil
}

func (mb *ModuleBuilder) BuildConstructor(g *types.GoType) error {
	err := mb.javaFile.WriteConstructorHeader(mb.className(g.PackageName), mb.constructorParams())
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteSuper(context)
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteCloseTag()
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteBlank(1)
	if err != nil {
		return err
	}
	return nil
}

func (mb *ModuleBuilder) className(packageName string) string {
	return mb.importedPackageName(packageName) + "Module"
}

func (mb *ModuleBuilder) importedPackageName(packageName string) string {
	return strings.Title(strings.ToLower(packageName))
}

func (mb *ModuleBuilder) BuildImports(g *types.GoType) error {
	err := mb.javaFile.WriteImport("com.facebook.react.bridge.ReactApplicationContext")
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteImport("com.facebook.react.bridge.Promise")
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteImport("com.facebook.react.bridge.ReactContextBaseJavaModule")
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteImport("com.facebook.react.bridge.ReactMethod")
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteImport("com.facebook.react.bridge.ReactMethod")
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteImport(mb.goImport(g.PackageName))
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteBlank(1)
	if err != nil {
		return err
	}
	return nil
}

func (mb *ModuleBuilder) goImport(packageName string) string {
	lowerPackage := strings.ToLower(packageName)
	return lowerPackage + "." + strings.Title(lowerPackage)
}

func (mb *ModuleBuilder) Close() error {
	return mb.javaFile.Close()
}

func (mb *ModuleBuilder) Create() error {
	return mb.javaFile.CreateFile()
}

func (mb *ModuleBuilder) constructorParams() map[string]string {
	params := make(map[string]string)
	params["ReactApplicationContext"] = context
	return params
}

func (mb *ModuleBuilder) BuildGetName(g *types.GoType) error {
	err := mb.javaFile.WriteAnnotation("Override")
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteMethodHeader("String", "getName", nil)
	if err != nil {
		return err
	}

	err = mb.javaFile.WriteReturn(mb.className(g.PackageName))
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteCloseTag()
	if err != nil {
		return err
	}
	return nil
}

func (mb *ModuleBuilder) BuildReactMethods(g *[]types.GoFunction, ret *[]types.GoParams, pkgName string) error {
	for i, val := range *g {
		loggers.Infof("BROKEN here with int %b, size: %b", i, len(*ret))
		err := mb.BuildReactMethod(&val, &(*ret)[i], pkgName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mb *ModuleBuilder) BuildReactMethod(g *types.GoFunction, ret *types.GoParams, pkgName string) error {
	err := mb.javaFile.WriteAnnotation("ReactMethod")

	if err != nil {
		return err
	}
	err = mb.buildMethodHeader(g, ret)
	if err != nil {
		return err
	}
	err = mb.buildReactMethodBody(g, ret, pkgName)
	if err != nil {
		return err
	}

	err = mb.javaFile.WriteCloseTag()
	if err != nil {
		return err
	}
	err = mb.javaFile.WriteBlank(1)
	if err != nil {
		return err
	}
	return nil
}

func (mb *ModuleBuilder) buildReactMethodBody(g *types.GoFunction, ret *types.GoParams, pkgName string) error {
	return mb.wrapTryCatch(func() error {
		return mb.methodMain(g, ret, pkgName)
	}, g, ret, "promise.reject(\"Error\")")
}

func (mb *ModuleBuilder) methodMain(g *types.GoFunction, ret *types.GoParams, pkgName string) error {
	methodCall := mb.importedPackageName(pkgName) + "." + g.Name + "()"
	if ret.T != "" {
		methodCall = types.GoToJava(ret.T) + " returnParam1 = " + methodCall
	}
	err := mb.javaFile.WriteMethodBody(methodCall)
	if err != nil {
		return err
	}
	if ret.T != "" {
		err = mb.javaFile.WriteMethodBody("promise.resolve(returnParam1)")
		if err != nil {
			return err
		}
	}
	return nil
}

func (mb *ModuleBuilder) wrapTryCatch(body func() error, g *types.GoFunction, ret *types.GoParams, catchMsg string) error {
	err := mb.javaFile.WriteTry()
	if err != nil {
		return err
	}
	err = body()
	if err != nil {
		return err
	}
	return mb.javaFile.WriteCatch("promise.reject(\"Error\", e)")
}

func (mb *ModuleBuilder) buildMethodHeader(g *types.GoFunction, ret *types.GoParams) error {
	paramMap := mb.ParamsToMap(g.Params)
	paramMap["Promise"] = "promise"
	return mb.javaFile.WriteMethodHeader("void", strings.ToLower(g.Name), paramMap)
}

func (mb *ModuleBuilder) ParamsToMap(params []types.GoParams) map[string]string {
	pmap := make(map[string]string)
	for _, k := range params {
		pmap[k.T] = k.Name
	}
	return pmap
}
