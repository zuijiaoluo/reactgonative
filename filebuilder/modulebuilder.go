package filebuilder

import (
	"path/filepath"
	"strings"

	"github.com/steve-winter/reactgonative/types"
)

var context = "reactContext"

// ModuleBuilder is the creator of each Classes boilerplate
type ModuleBuilder struct {
	javaFile *JavaFile
}

// NewModuleBuilder returns a new ModuleBuilder containing a JavaFile.
// The file is not opened or created at this point.
func NewModuleBuilder(name string, root string) ModuleBuilder {
	return ModuleBuilder{
		javaFile: NewJavaFile(name, root),
	}
}

func (mb *ModuleBuilder) createPackageName(name string, root string) string {
	var line string
	if root == "" {
		line = "bridge." + name
	} else {
		line = root + ".bridge." + name
	}
	return line
}

func (mb *ModuleBuilder) buildFileName(pkgName string, pkgRoot string) string {
	packageNameString := strings.Replace(
		mb.createPackageName(pkgName, pkgRoot), ".", "/", -1)
	fileName := filepath.Join(mb.javaFile.fileName,
		packageNameString)
	// dir, _ := filepath.Split(fileName)
	fileName = filepath.Join(fileName, mb.className(pkgName)+".java")
	return fileName
}

//BuildModule generates the Java class with features in g.
//Returns the className created, or an error if a write fails
func (mb *ModuleBuilder) BuildModule(g *types.GoType) (string, error) {
	fileName := mb.buildFileName(g.PackageName, mb.javaFile.packageRoot)
	mb.javaFile.setFileName(fileName)
	err := mb.create()
	if err != nil {
		return "", err
	}
	err = mb.javaFile.writePackageLine(mb.createPackageName(g.PackageName, mb.javaFile.packageRoot))
	if err != nil {
		return "", err
	}
	err = mb.javaFile.writeBlank(1)
	if err != nil {
		return "", err
	}
	err = mb.buildImports(g)
	if err != nil {
		return "", err
	}

	err = mb.javaFile.writeClassHeader(mb.className(g.PackageName),
		"ReactContextBaseJavaModule", "")
	if err != nil {
		return "", err
	}
	err = mb.javaFile.writeBlank(1)
	if err != nil {
		return "", err
	}
	err = mb.buildConstructor(g)
	if err != nil {
		return "", err
	}

	err = mb.buildGetName(g)
	if err != nil {
		return "", err
	}
	err = mb.javaFile.writeBlank(1)
	if err != nil {
		return "", err
	}

	err = mb.buildReactMethods(&g.Functions, &g.Returns, g.PackageName)
	if err != nil {
		return "", err
	}
	err = mb.javaFile.writeCloseTag()
	if err != nil {
		return "", err
	}
	return mb.className(g.PackageName), nil
}

func (mb *ModuleBuilder) buildConstructor(g *types.GoType) error {
	err := mb.javaFile.writeConstructorHeader(mb.className(g.PackageName), mb.constructorParams())
	if err != nil {
		return err
	}
	err = mb.javaFile.writeSuper(context)
	if err != nil {
		return err
	}
	err = mb.javaFile.writeCloseTag()
	if err != nil {
		return err
	}
	err = mb.javaFile.writeBlank(1)
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

func (mb *ModuleBuilder) buildImports(g *types.GoType) error {
	err := mb.javaFile.writeImport("com.facebook.react.bridge.ReactApplicationContext")
	if err != nil {
		return err
	}
	err = mb.javaFile.writeImport("com.facebook.react.bridge.Promise")
	if err != nil {
		return err
	}
	err = mb.javaFile.writeImport("com.facebook.react.bridge.ReactContextBaseJavaModule")
	if err != nil {
		return err
	}
	err = mb.javaFile.writeImport("com.facebook.react.bridge.ReactMethod")
	if err != nil {
		return err
	}
	err = mb.javaFile.writeImport("com.facebook.react.bridge.ReactMethod")
	if err != nil {
		return err
	}
	err = mb.javaFile.writeImport(mb.goImport(g.PackageName))
	if err != nil {
		return err
	}
	err = mb.javaFile.writeBlank(1)
	if err != nil {
		return err
	}
	return nil
}

func (mb *ModuleBuilder) goImport(packageName string) string {
	lowerPackage := strings.ToLower(packageName)
	return lowerPackage + "." + strings.Title(lowerPackage)
}

// Close will close the internal javaFile
func (mb *ModuleBuilder) Close() error {
	return mb.javaFile.close()
}

func (mb *ModuleBuilder) create() error {
	return mb.javaFile.createFile()
}

func (mb *ModuleBuilder) constructorParams() []types.GoParams {
	params := make([]types.GoParams, 0)
	params = append(params, types.GoParams{Name: context, T: "ReactApplicationContext"})
	return params
}

func (mb *ModuleBuilder) buildGetName(g *types.GoType) error {
	err := mb.javaFile.writeAnnotation("Override")
	if err != nil {
		return err
	}
	err = mb.javaFile.writeMethodHeader("String", "getName", nil)
	if err != nil {
		return err
	}

	err = mb.javaFile.writeReturnStatic(mb.className(g.PackageName))
	if err != nil {
		return err
	}
	err = mb.javaFile.writeCloseTag()
	if err != nil {
		return err
	}
	return nil
}

func (mb *ModuleBuilder) buildReactMethods(g *[]types.GoFunction, ret *[]types.GoParams, pkgName string) error {
	for i, val := range *g {
		err := mb.buildReactMethod(&val, &(*ret)[i], pkgName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mb *ModuleBuilder) buildReactMethod(g *types.GoFunction, ret *types.GoParams, pkgName string) error {
	err := mb.javaFile.writeAnnotation("ReactMethod")

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

	err = mb.javaFile.writeCloseTag()
	if err != nil {
		return err
	}
	err = mb.javaFile.writeBlank(1)
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

	methodCall := mb.importedPackageName(pkgName) + "." + strings.ToLower(g.Name) + "(" + mb.buildMethodCallParams(&g.Params) + ")"
	if ret.T != "" {
		methodCall = types.GoToJava(ret.T) + " returnParam1 = " + methodCall
	}
	err := mb.javaFile.writeMethodBody(methodCall)
	if err != nil {
		return err
	}
	if ret.T != "" {
		err = mb.javaFile.writeMethodBody("promise.resolve(returnParam1)")
		if err != nil {
			return err
		}
	}
	return nil
}

func (mb *ModuleBuilder) buildMethodCallParams(g *[]types.GoParams) string {
	paramsMap := mb.paramsToMap(*g)
	resp := ""
	for _, val := range paramsMap {
		if len(resp) != 0 {
			resp = resp + ", "
		}
		resp = resp + val.Name
	}
	return resp
}

func (mb *ModuleBuilder) wrapTryCatch(body func() error, g *types.GoFunction, ret *types.GoParams, catchMsg string) error {
	err := mb.javaFile.writeTry()
	if err != nil {
		return err
	}
	err = body()
	if err != nil {
		return err
	}
	return mb.javaFile.writeCatch("promise.reject(\"Error\", e)")
}

func (mb *ModuleBuilder) buildMethodHeader(g *types.GoFunction, ret *types.GoParams) error {
	params := make([]types.GoParams, 0)
	params = append(params, *ret)
	params = append(params, types.GoParams{Name: "promise", T: "Promise"})
	return mb.javaFile.writeMethodHeader("void", strings.ToLower(g.Name), params)
}

func (mb *ModuleBuilder) paramsToMap(params []types.GoParams) []types.GoParams {
	// count := 0
	// for v := range params {
	// 	for _, k := range v {
	// 		pmap[k.Name] = k.T
	// 	}
	// }
	return params
}
