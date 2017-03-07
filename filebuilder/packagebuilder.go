package filebuilder

import (
	"path/filepath"
	"strings"

	"github.com/steve-winter/loggers"
)

type PackageBuilder struct {
	javaFile *JavaFile
}

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

func (pb *PackageBuilder) BuildPackage(typeString string, packageName string) error {
	fileName := pb.buildFileName(packageName, pb.javaFile.packageRoot)
	pb.javaFile.SetFileName(fileName)
	loggers.Infof("Filename is: %s", fileName)
	err := pb.Create()
	if err != nil {
		return err
	}
	err = pb.javaFile.WritePackageLine(pb.createPackageName(packageName, pb.javaFile.packageRoot))
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteBlank(1)
	if err != nil {
		return err
	}
	err = pb.BuildImports()
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteBlank(1)
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteClassHeader(pb.className(packageName),
		"", "ReactPackage")
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteBlank(1)
	if err != nil {
		return err
	}
	err = pb.BuildNativeModulesMethod(packageName)
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteBlank(1)
	if err != nil {
		return err
	}
	err = pb.BuildCreateJSModulesMethod()
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteBlank(1)
	if err != nil {
		return err
	}
	err = pb.BuildCreateViewManagersMethod()
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteBlank(1)
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteCloseTag()
	if err != nil {
		return err
	}

	return nil
}

func (pb *PackageBuilder) BuildImports() error {
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
			err := pb.javaFile.WriteImport(val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (pb *PackageBuilder) BuildNativeModulesMethod(packageName string) error {
	params := make(map[string]string)
	params["ReactApplicationContext"] = context
	err := pb.javaFile.WriteAnnotation("Override")
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteMethodHeader("List<NativeModule>", "createNativeModules", params)
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteMethodBody("List<NativeModule> modules = new ArrayList<>()")
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteMethodBody("modules.add(new " + pb.moduleName(packageName) + "(" + context + "))")
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteReturnDynamic("modules")
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteCloseTag()
	if err != nil {
		return err
	}
	return nil
}

func (pb *PackageBuilder) BuildCreateJSModulesMethod() error {
	err := pb.javaFile.WriteAnnotation("Override")
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteMethodHeader("List <Class<? extends JavaScriptModule>>", "createJSModules", nil)
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteReturnDynamic("Collections.emptyList()")
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteCloseTag()
	if err != nil {
		return err
	}
	return nil
}

func (pb *PackageBuilder) BuildCreateViewManagersMethod() error {
	params := make(map[string]string)
	params["ReactApplicationContext"] = context
	err := pb.javaFile.WriteAnnotation("Override")
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteMethodHeader("List<ViewManager>", "createViewManagers", params)
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteReturnDynamic("Collections.emptyList()")
	if err != nil {
		return err
	}
	err = pb.javaFile.WriteCloseTag()
	if err != nil {
		return err
	}
	return nil
}

func (mb *PackageBuilder) Close() error {
	return mb.javaFile.Close()
}

func (mb *PackageBuilder) Create() error {
	return mb.javaFile.CreateFile()
}

func (mb *PackageBuilder) className(packageName string) string {
	return mb.importedPackageName(packageName) + "Package"
}

func (mb *PackageBuilder) moduleName(packageName string) string {
	return mb.importedPackageName(packageName) + "Module"
}

func (mb *PackageBuilder) importedPackageName(packageName string) string {
	return strings.Title(strings.ToLower(packageName))
}
