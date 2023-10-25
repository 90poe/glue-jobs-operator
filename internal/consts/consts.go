package consts

const (
	// Error types for GlueJob
	// 2 types: RevocerableError and UnrecoverableError
	RecoverableError   = "RecoverableError"
	UnrecoverableError = "UnrecoverableError"
	SuccessReconcile   = "Success"
	// GlueJob status Type
	StatusReady    = "Ready"
	StatusNotReady = "NotReady"
)
