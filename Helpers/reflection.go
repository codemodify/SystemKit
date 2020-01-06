package Helpers

import (
	"fmt"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
)

func getFuncName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	var functionName = runtime.FuncForPC(pc).Name()

	functionName = strings.Replace(functionName, strToReplace, strToReplaceWith, 1)
	functionName = strings.Replace(functionName, "(", "", 1)
	functionName = strings.Replace(functionName, ")", "", 1)
	functionName = strings.Replace(functionName, "*", "", 1)

	return functionName + "()"
}

// GetLineNumber -
func GetLineNumber() int {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Line
}

// GetThisFuncName -
func GetThisFuncName() string {
	return getFuncName(2)
}

// GetThisFuncNameWithLine -
func GetThisFuncNameWithLine() string {
	return fmt.Sprintf("%d : %s", getFuncName(2), GetLineNumber())
}

// GetStackTrace -
func GetStackTrace2(prefixToRemove string) string {

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var stack = ""

	frame, more := frames.Next()
	for more {
		// stack += newLineSeparator
		funcName := strings.ReplaceAll(frame.Function, prefixToRemove, "")
		stack += funcName + "()/" + stack

		frame, more = frames.Next()
	}

	if len(stack) > 0 {
		stack = stack[0 : len(stack)-1]
	}

	return stack
}

// GetStackTrace -
func GetStackTrace(prefixesToRemove []string) string {
	// debug.PrintStack()

	prefixesToRemove = append(prefixesToRemove, "created by ")
	prefixesToRemove = append(prefixesToRemove, ".()")

	stackTrace := debug.Stack()
	lines := strings.Split(string(stackTrace), "\n")

	stack := ""
	for _, funcName := range lines {
		if strings.Index(funcName, "goroutine") != -1 ||
			strings.Index(funcName, "runtime/debug") != -1 ||
			strings.Index(funcName, ".go:") != -1 ||
			strings.Index(funcName, "GetStackTrace") != -1 ||
			len(strings.TrimSpace(funcName)) <= 0 {
			continue
		}

		regExp := regexp.MustCompile(`\((.*?)\)`)
		regExpResult := regExp.FindAllStringSubmatch(funcName, -1)
		for _, arr := range regExpResult {
			if len(arr) > 1 {
				toRemove := arr[1]
				if strings.Index(toRemove, "0x") != -1 {
					funcName = strings.Replace(funcName, toRemove, "", -1)
				}
			}
		}

		for _, prefixeToRemove := range prefixesToRemove {
			funcName = strings.Replace(funcName, prefixeToRemove, "", -1)
		}

		if strings.Index(funcName, "(") == -1 {
			funcName += "()"
		}

		if len(stack) > 0 {
			stack = stack + " | " + funcName
		} else {
			stack = funcName
		}
	}

	return stack
}
