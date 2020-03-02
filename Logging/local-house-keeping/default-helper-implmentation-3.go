package housekeeping

func (thisRef defaultHelperImplmentation) LogPanic(message string) {
	thisRef.LogPanicWithTagAndLevel("", 0, message)
}
func (thisRef defaultHelperImplmentation) LogFatal(message string) {
	thisRef.LogFatalWithTagAndLevel("", 0, message)
}
func (thisRef defaultHelperImplmentation) LogError(message string) {
	thisRef.LogErrorWithTagAndLevel("", 0, message)
}
func (thisRef defaultHelperImplmentation) LogWarning(message string) {
	thisRef.LogWarningWithTagAndLevel("", 0, message)
}
func (thisRef defaultHelperImplmentation) LogInfo(message string) {
	thisRef.LogInfoWithTagAndLevel("", 0, message)
}
func (thisRef defaultHelperImplmentation) LogSuccess(message string) {
	thisRef.LogSuccessWithTagAndLevel("", 0, message)
}
func (thisRef defaultHelperImplmentation) LogDebug(message string) {
	thisRef.LogDebugWithTagAndLevel("", 0, message)
}
func (thisRef defaultHelperImplmentation) LogTrace(message string) {
	thisRef.LogTraceWithTagAndLevel("", 0, message)
}
