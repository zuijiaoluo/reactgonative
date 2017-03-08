package filebuilder

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/steve-winter/reactgonative/types"
)

type JavaFile struct {
	f            *os.File
	fileName     string
	packageRoot  string
	depth        int
	shouldIndent bool
}

func NewJavaFile(name string, root string) (javaFile *JavaFile) {
	return &JavaFile{
		fileName:     name,
		packageRoot:  root,
		shouldIndent: false,
	}
}

func (jf *JavaFile) setFileName(name string) error {
	if jf.f != nil {
		return errors.New("File already open")
	}
	jf.fileName = name
	return nil
}

func (jf *JavaFile) createFile() error {
	dir, _ := filepath.Split(jf.fileName)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	f, err := os.Create(jf.fileName)
	if err != nil {
		return err
	}
	jf.f = f
	return nil
}

func (jf *JavaFile) writePackageLine(packageName string) error {
	err := jf.writeLineFlat("package " + packageName + ";")
	return err
}

func (jf *JavaFile) writeImport(importLine string) error {
	return jf.writeLineFlat("import " + importLine + ";")
}

func (jf *JavaFile) writeClassHeader(className string, extendsName string, implementsName string) error {
	line := "public class " + className + jf.extends(extendsName) + jf.implements(implementsName)
	return jf.writeLine(line + " {")
}

func (jf *JavaFile) writeConstructorHeader(className string, params map[string]string) error {
	return jf.writeLine(jf.constructorHeader(className, params))
}

func (jf *JavaFile) constructorHeader(className string, params map[string]string) string {
	line := "public " + className + "(" + jf.methodParams(params) + ") {"
	return line
}

func (jf *JavaFile) writeSuper(param string) error {
	return jf.writeLine("super(" + param + ");")
}

func (jf *JavaFile) writeAnnotation(annot string) error {
	return jf.writeLineN("@" + annot)
}

func (jf *JavaFile) writeMethodHeader(returnType string, methodName string, params map[string]string) error {
	return jf.writeLine("public " + returnType + " " + methodName + "(" + jf.methodParams(params) + ") {")
}

func (jf *JavaFile) writeMethodBody(body string) error {
	return jf.writeLineN(body + ";")
}

func (jf *JavaFile) writeCloseTag() error {
	return jf.writeLine("}")
}

func (jf *JavaFile) methodParams(params map[string]string) string {
	line := ""
	for j, k := range params {
		if len(line) > 0 {
			line = line + ", "
		}
		line = line + types.GoToJava(j) + " " + k
	}
	return line
}

func (jf *JavaFile) extends(extendsName string) string {
	return jf.classModifier("extends", extendsName)
}

func (jf *JavaFile) implements(implementsName string) string {
	return jf.classModifier("implements", implementsName)
}

func (jf *JavaFile) classModifier(modifier string, modifierName string) string {
	if len(modifierName) > 0 {
		return " " + modifier + " " + modifierName
	}
	return ""
}

func (jf *JavaFile) writeTry() error {
	return jf.writeLine("try {")
}

func (jf *JavaFile) writeCatch(msg string) error {
	err := jf.writeLine("} catch(Exception e) {")
	if err != nil {
		return err
	}
	err = jf.writeLine(msg + ";")
	if err != nil {
		return err
	}
	return jf.writeLineN("}")
}

func (jf *JavaFile) writeReturnDynamic(ret string) error {
	return jf.writeLine("return " + ret + ";")
}

func (jf *JavaFile) writeReturnStatic(ret string) error {
	return jf.writeLine("return \"" + ret + "\";")

}

func (jf *JavaFile) writeLine(line string) error {
	indent := ""
	jf.shouldStep(line)
	for i := 0; i < jf.depth; i++ {
		indent = indent + "\t"
	}
	err := jf.writeLineFlat(indent + line)
	if !strings.EqualFold(line, "}") {
		jf.shouldIndent = true
	}
	return err
}

func (jf *JavaFile) writeLineN(line string) error {
	indent := ""
	jf.shouldStep(line)
	for i := 0; i < jf.depth; i++ {
		indent = indent + "\t"
	}
	jf.shouldIndent = false
	return jf.writeLineFlat(indent + line)
}

func (jf *JavaFile) shouldStep(line string) {
	if strings.EqualFold(line, "}") {
		jf.depth = jf.depth - 1
		jf.shouldIndent = false
		return
	}
	if strings.Contains(line, "}") {
		jf.depth = jf.depth - 1
		// oldShouldIndent := jf.shouldIndent
		jf.shouldIndent = strings.Contains(line, "{")
		// if !jf.shouldIndent {
		// 	jf.depth = jf.depth + 1
		// } else {
		// 	jf.depth = jf.depth - 1
		// }
		// if oldShouldIndent {
		// 	jf.shouldIndent = oldShouldIndent
		// }
		return
	}
	if jf.shouldIndent {
		jf.depth = jf.depth + 1
		return
	}
}

func (jf *JavaFile) writeLineFlat(line string) error {
	_, err := jf.f.WriteString(line + "\n")
	if err == nil {
		err = jf.f.Sync()
	}
	return err
}

func (jf *JavaFile) writeBlank(num int) error {
	newLine := ""
	for i := 0; i < num; i++ {
		newLine = newLine + "\n"
	}
	_, err := jf.f.WriteString(newLine)
	return err
}

func (jf *JavaFile) close() error {

	return jf.f.Close()
}
