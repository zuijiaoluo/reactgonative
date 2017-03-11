package filebuilder

import (
	"path/filepath"
	"strings"

	"github.com/steve-winter/reactgonative/types"
)

// PackageBuilder is the creator of each Packages boilerplate
type PackageBuilder struct {
	javaFile *JavaFile
}

// NewPackageBuilder returns a new PackageBuilder containing a JavaFile.
// The file is not opened or created at this point.
func NewPackageBuilder(name string, root string) PackageBuilder {
	return PackageBuilder{
		javaFile: NewJavaFile(name, root),
	}
}

func (pb *PackageBuilder) createPackageName(name string, root string) string {
	var line string
	if root == "" {
		line = "bridge." + name
	} else {
		line = root + ".bridge." + name
	}
	return line
}

func (pb *PackageBuilder) buildFileName(pkgName string, pkgRoot string) string {
	packageNameString := strings.Replace(
		pb.createPackageName(pkgName, pkgRoot), ".", "/", -1)
	fileName := filepath.Join(pb.javaFile.fileName,
		packageNameString)
	// dir, _ := filepath.Split(fileName)
	fileName = filepath.Join(fileName, pb.className(pkgName)+".java")

	return fileName
}

// BuildPackage generates the package boilerplate for the inputted packageName.
// An error is returned if any write fail
func (pb *PackageBuilder) BuildPackage(packageName string) error {
	fileName := pb.buildFileName(packageName, pb.javaFile.packageRoot)
	pb.javaFile.setFileName(fileName)
	err := pb.create()
	if err != nil {
		return err
	}
	err = pb.javaFile.writePackageLine(pb.createPackageName(packageName, pb.javaFile.packageRoot))
	if err != nil {
		return err
	}
	err = pb.javaFile.writeBlank(1)
	if err != nil {
		return err
	}
	err = pb.buildImports()
	if err != nil {
		return err
	}
	err = pb.javaFile.writeBlank(1)
	if err != nil {
		return err
	}
	err = pb.javaFile.writeClassHeader(pb.className(packageName),
		"", "ReactPackage")
	if err != nil {
		return err
	}
	err = pb.javaFile.writeBlank(1)
	if err != nil {
		return err
	}
	err = pb.buildNativeModulesMethod(packageName)
	if err != nil {
		return err
	}
	err = pb.javaFile.writeBlank(1)
	if err != nil {
		return err
	}
	err = pb.buildCreateJSModulesMethod()
	if err != nil {
		return err
	}
	err = pb.javaFile.writeBlank(1)
	if err != nil {
		return err
	}
	err = pb.buildCreateViewManagersMethod()
	if err != nil {
		return err
	}
	err = pb.javaFile.writeBlank(1)
	if err != nil {
		return err
	}
	err = pb.javaFile.writeCloseTag()
	if err != nil {
		return err
	}

	return nil
}

func (pb *PackageBuilder) buildImports() error {
	imports := [8]string{
		"com.facebook.react.ReactPackage",
		"com.facebook.react.bridge.JavaScriptModule",
		"com.facebook.react.bridge.NativeModule",
		"com.facebook.react.bridge.ReactApplicationContext",
		"com.facebook.react.uimanager.ViewManager",
		"java.util.ArrayList",
		"java.util.Collections",
		"java.util.List",
	}
	for _, val := range imports {
		if !strings.EqualFold(val, "") {
			err := pb.javaFile.writeImport(val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (pb *PackageBuilder) buildNativeModulesMethod(packageName string) error {
	params := make([]types.GoParams, 0)
	params = append(params, types.GoParams{Name: context, T: "ReactApplicationContext"})
	err := pb.javaFile.writeAnnotation("Override")
	if err != nil {
		return err
	}
	err = pb.javaFile.writeMethodHeader("List<NativeModule>", "createNativeModules", params)
	if err != nil {
		return err
	}
	err = pb.javaFile.writeMethodBody("List<NativeModule> modules = new ArrayList<>()")
	if err != nil {
		return err
	}
	err = pb.javaFile.writeMethodBody("modules.add(new " + pb.moduleName(packageName) + "(" + context + "))")
	if err != nil {
		return err
	}
	err = pb.javaFile.writeReturnDynamic("modules")
	if err != nil {
		return err
	}
	err = pb.javaFile.writeCloseTag()
	if err != nil {
		return err
	}
	return nil
}

func (pb *PackageBuilder) buildCreateJSModulesMethod() error {
	err := pb.javaFile.writeAnnotation("Override")
	if err != nil {
		return err
	}
	err = pb.javaFile.writeMethodHeader("List <Class<? extends JavaScriptModule>>", "createJSModules", nil)
	if err != nil {
		return err
	}
	err = pb.javaFile.writeReturnDynamic("Collections.emptyList()")
	if err != nil {
		return err
	}
	err = pb.javaFile.writeCloseTag()
	if err != nil {
		return err
	}
	return nil
}

func (pb *PackageBuilder) buildCreateViewManagersMethod() error {
	params := make([]types.GoParams, 0)
	params = append(params, types.GoParams{Name: context, T: "ReactApplicationContext"})
	err := pb.javaFile.writeAnnotation("Override")
	if err != nil {
		return err
	}
	err = pb.javaFile.writeMethodHeader("List<ViewManager>", "createViewManagers", params)
	if err != nil {
		return err
	}
	err = pb.javaFile.writeReturnDynamic("Collections.emptyList()")
	if err != nil {
		return err
	}
	err = pb.javaFile.writeCloseTag()
	if err != nil {
		return err
	}
	return nil
}

// Close will close the internal JavaFile
func (pb *PackageBuilder) Close() error {
	return pb.javaFile.close()
}

func (pb *PackageBuilder) create() error {
	return pb.javaFile.createFile()
}

func (pb *PackageBuilder) className(packageName string) string {
	return pb.importedPackageName(packageName) + "Package"
}

func (pb *PackageBuilder) moduleName(packageName string) string {
	return pb.importedPackageName(packageName) + "Module"
}

func (pb *PackageBuilder) importedPackageName(packageName string) string {
	return strings.Title(strings.ToLower(packageName))
}
