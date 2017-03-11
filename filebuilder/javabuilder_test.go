package filebuilder

import (
	"bufio"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewJavaFile(t *testing.T) {
	Convey("Given a creating a new JavaFile", t, func() {
		Convey("When creating a new file with given fields", func() {
			jf := NewJavaFile("testFileName", "testFileRoot")
			Convey("Then the filename is testFileName", func() {
				So(jf.fileName, ShouldEqual, "testFileName")
			})
			Convey("And the fileRoot is fileRoot", func() {
				So(jf.packageRoot, ShouldEqual, "testFileRoot")
			})
			Convey("And the depth is 0", func() {
				So(jf.depth, ShouldEqual, 0)
			})
			Convey("And shouldIndent is false", func() {
				So(jf.shouldIndent, ShouldEqual, false)
			})
			Convey("And file is nil", func() {
				So(jf.f, ShouldEqual, nil)
			})
		})
	})
}

func TestSetFileName(t *testing.T) {
	Convey("Given javafile created", t, func() {
		jf := NewJavaFile("testFileName", "testFileRoot")
		Convey("When filename is set to helloWorld", func() {
			err := jf.setFileName("helloWorld")
			Convey("Then the filename changes to helloWorld", func() {
				So(jf.fileName, ShouldEqual, "helloWorld")
			})
			Convey("And the error should be nil", func() {
				So(err, ShouldEqual, nil)
			})
		})
	})
	Convey("Given javafile created and file opened", t, func() {
		jf := NewJavaFile("/tmp/reactgonative/testfile_setFileName", "testFileRoot")
		jf.createFile()
		Convey("When filename is set to helloWorld", func() {
			err := jf.setFileName("helloWorld")
			Convey("Then the filename should not change to helloWorld", func() {
				So(jf.fileName, ShouldNotEqual, "helloWorld")
			})
			Convey("And the error should be File is already open", func() {
				So(err.Error(), ShouldEqual, "File already open")
			})
		})
	})
}

func TestCreateFile(t *testing.T) {
	Convey("Given javafile object created", t, func() {
		jf := NewJavaFile("/tmp/reactgonative/testfile_createFile1", "testFileRoot")
		Convey("When a file is created", func() {
			err := jf.createFile()
			Convey("Then no error is generated", func() {
				So(err, ShouldEqual, nil)
			})
			Convey("And the file is set", func() {
				So(jf.f, ShouldNotEqual, nil)
			})
			Convey("And the filename of file is same as JavaFile", func() {
				So(jf.f.Name(), ShouldEqual, jf.fileName)
			})
			Convey("When the same file is created", func() {
				sf := NewJavaFile("/tmp/reactgonative/testfile_createFile1", "testFileRoot")
				err2 := sf.createFile()
				Convey("Then no error is generated", func() {
					So(err2, ShouldEqual, nil)
				})
			})
		})
	})
	Convey("Given file created", t, func() {
		cf := NewJavaFile("/tmp/reactgonative/testfile_createFileCreated", "testFileRoot")
		cf.createFile()
		Convey("When a new file is created using same directory", func() {
			jf := NewJavaFile("/tmp/reactgonative/testfile_createFileCreated/createFile2", "testFileRoot")
			err := jf.createFile()
			Convey("Then the folder generation should generate an error", func() {
				So(err, ShouldNotEqual, nil)
			})
		})
	})
	Convey("Given folder created", t, func() {
		cf := NewJavaFile("/tmp/reactgonative/testfile_createFileCreated_folder/somefile", "testFileRoot")
		cf.createFile()
		Convey("When a new file is created using filename of existing folder", func() {
			jf := NewJavaFile("/tmp/reactgonative/testfile_createFileCreated_folder", "testFileRoot")
			err := jf.createFile()
			Convey("Then the file generation should generate an error", func() {
				So(err.Error(), ShouldEndWith, "is a directory")
			})
		})
	})
}

func TestWritePackageLine(t *testing.T) {
	Convey("Given Javafile object created", t, func() {
		cf := NewJavaFile("/tmp/reactgonative/testfile_writePackageLine1", "testFileRoot")
		Convey("When the file is not open and write attempted", func() {
			err := cf.writePackageLine("packageName")
			Convey("Then an error is generated", func() {
				So(err.Error(), ShouldEqual, "invalid argument")
			})
		})
		Convey("When the file is open and write attempted", func() {
			cf.createFile()
			err := cf.writePackageLine("packageName")
			Convey("Then no error is generated", func() {
				So(err, ShouldEqual, nil)
			})
			Convey("And the end line of the file is packageName", func() {
				So(readLastLines("/tmp/reactgonative/testfile_writePackageLine1", 1)[0],
					ShouldEqual, "package packageName;")
			})
		})
	})
}

func TestWriteImport(t *testing.T) {
	Convey("Given Javafile object created", t, func() {
		cf := NewJavaFile("/tmp/reactgonative/testfile_writeImport1", "testFileRoot")
		Convey("When the file is open and write attempted", func() {
			cf.createFile()
			err := cf.writeImport("com.lemonade.pink")
			Convey("Then no error is generated", func() {
				So(err, ShouldEqual, nil)
			})
			Convey("And the end line of the file is packageName", func() {
				So(readLastLines("/tmp/reactgonative/testfile_writeImport1", 1)[0],
					ShouldEqual, "import com.lemonade.pink;")
			})
		})
	})
}

func TestWriteClassHeader(t *testing.T) {
	Convey("Given Javafile object created", t, func() {
		cf := NewJavaFile("/tmp/reactgonative/testfile_writeClassHeader1", "testFileRoot")
		Convey("When the file is open and writing class header attempted", func() {
			cf.createFile()
			err := cf.writeClassHeader("MyClassName", "ExtendsName", "InterfaceName")
			Convey("Then no error is generated", func() {
				So(err, ShouldEqual, nil)
			})
			Convey("And the end line of the file is packageName", func() {
				So(readLastLines("/tmp/reactgonative/testfile_writeClassHeader1", 1)[0],
					ShouldEqual, "public class MyClassName extends ExtendsName implements InterfaceName {")
			})
		})
	})
}

func readLastLines(fileName string, lineCount int) []string {
	fileHandle, _ := os.Open(fileName)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)
	lines := make([]string, lineCount)
	count := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lines[count%lineCount] = line
		count++
	}
	return lines
}
