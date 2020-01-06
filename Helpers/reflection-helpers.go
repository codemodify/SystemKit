package Helpers

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
)

var strToReplace = "github.com/"
var strToReplaceWith = ""

var mainPackagePath string

// InitMainPackagePath - to be called in `main()`
func InitMainPackagePath() {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		mainPackagePath = path.Dir(filename) + "/"
		fmt.Println(mainPackagePath)

		mainPackagePath = strings.Replace(mainPackagePath, strToReplace, strToReplaceWith, 1)
	}
}

// GetPathForThisPackage -
func GetPathForThisPackage() string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return path.Dir(filename) + "/"
	}

	return "./"
}

// GetStructName - returns struct name
func GetStructName(i interface{}) string {
	var t = reflect.TypeOf(i)

	var result string

	if t.Kind() == reflect.Ptr {
		result = t.Elem().Name()
	} else {
		result = t.Name()
	}

	return strings.Replace(result, strToReplace, strToReplaceWith, 1)
}

// GetPackageName - returns package name
func GetPackageName(i interface{}) string {
	var t = reflect.TypeOf(i)

	var result string

	if t.Kind() == reflect.Ptr {
		result = t.Elem().PkgPath()
	} else {
		result = t.PkgPath()
	}

	return strings.Replace(result, strToReplace, strToReplaceWith, 1)
}

// GetPackageNameWithStruct - returns package name with struct
func GetPackageNameWithStruct(i interface{}) string {
	var t = reflect.TypeOf(i)

	var result string

	if t.Kind() == reflect.Ptr {
		result = t.Elem().PkgPath() + "/" + t.Elem().Name()
	} else {
		result = t.PkgPath() + "/" + t.Name()
	}

	return strings.Replace(result, strToReplace, strToReplaceWith, 1)
}

// GetCallingFuncName - returns package + function name at runtime
func GetCallingFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	var functionName = runtime.FuncForPC(pc).Name()

	functionName = strings.Replace(functionName, strToReplace, strToReplaceWith, 1)
	functionName = strings.Replace(functionName, "(", "", 1)
	functionName = strings.Replace(functionName, ")", "", 1)
	functionName = strings.Replace(functionName, "*", "", 1)

	return functionName + "()"
}

// GetCallStackWithFileAndLineNumber - traces a call with line number
func GetCallStackWithFileAndLineNumber() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var logLine = ""

	frame, more := frames.Next()
	for more {
		frame.File = strings.Replace(frame.File, mainPackagePath, "", 1)
		// logLine += fmt.Sprintf("%s,:%d %s", frame.File, frame.Line, frame.Function)
		logLine += fmt.Sprintf("%s | ", frame.Function)

		frame, more = frames.Next()
	}

	return logLine
}

// CallStack - returns function name at runtime
func CallStack() string {
	pc, _, _, _ := runtime.Caller(1)
	var functionName = runtime.FuncForPC(pc).Name()

	return strings.Replace(functionName, "_"+mainPackagePath, "", 1)
}

// PackageName - returns function name at runtime
func PackageName(i interface{}) string {
	var packagePath = reflect.TypeOf(i).PkgPath()
	runes := []rune(packagePath)
	return string(runes[1:len(packagePath)])
}

// Trace - traces a call with line number
func Trace() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fmt.Printf("%s,:%d %s\n", frame.File, frame.Line, frame.Function)
}
