package filebuilder

import (
	"path/filepath"
	"strings"
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
	pb.javaFile.setFileName(fileName)
	err := pb.Create()
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
	params := make(map[string]string)
	params["ReactApplicationContext"] = context
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
	params := make(map[string]string)
	params["ReactApplicationContext"] = context
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

func (mb *PackageBuilder) Close() error {
	return mb.javaFile.close()
}

func (mb *PackageBuilder) Create() error {
	return mb.javaFile.createFile()
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
