package housekeeping

import (
	"encoding/json"
	"fmt"

	loggingC "github.com/codemodify/SystemKit/Logging/Contracts"
)

func (thisRef defaultHelperImplmentation) LogPanicWithFields(fields loggingC.Fields) {
	var data, err = json.Marshal(fields)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}

	thisRef.LogPanicWithTagAndLevel("", 0, string(data))
}
func (thisRef defaultHelperImplmentation) LogFatalWithFields(fields loggingC.Fields) {
	var data, err = json.Marshal(fields)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}

	thisRef.LogFatalWithTagAndLevel("", 0, string(data))
}
func (thisRef defaultHelperImplmentation) LogErrorWithFields(fields loggingC.Fields) {
	var data, err = json.Marshal(fields)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}

	thisRef.LogErrorWithTagAndLevel("", 0, string(data))
}
func (thisRef defaultHelperImplmentation) LogWarningWithFields(fields loggingC.Fields) {
	var data, err = json.Marshal(fields)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}

	thisRef.LogWarningWithTagAndLevel("", 0, string(data))
}
func (thisRef defaultHelperImplmentation) LogInfoWithFields(fields loggingC.Fields) {
	var data, err = json.Marshal(fields)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}

	thisRef.LogInfoWithTagAndLevel("", 0, string(data))
}
func (thisRef defaultHelperImplmentation) LogDebugWithFields(fields loggingC.Fields) {
	var data, err = json.Marshal(fields)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}

	thisRef.LogDebugWithTagAndLevel("", 0, string(data))
}
func (thisRef defaultHelperImplmentation) LogTraceWithFields(fields loggingC.Fields) {
	var data, err = json.Marshal(fields)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}

	thisRef.LogTraceWithTagAndLevel("", 0, string(data))
}
