package codes

import "fmt"

const (
	OK = 0
	ConfigErrorCode = 101
	NetworkNotSupportedErrorCode = 201
	ClientMsgErrorCode = 301
	ClientCertFail = 401
)

// errorcode type
const (
	FrameworkError = 1
	BusinuessError = 2
)


// framework error
var (
	ConfigError = NewFrameworkError(ConfigErrorCode,"config error")
	NetworkNotSupportedError = NewFrameworkError(NetworkNotSupportedErrorCode,"network type not supported")
	ClientCertFailError = NewFrameworkError(ClientCertFail, "client cert fail")
)


// Error defines all errors in the framework
type Error struct {
	Code int
	Type int
	Message string
}

const (
	Success = "success"
)

func (e *Error) Error() string {
	if e == nil {
		return Success
	}
	if e.Type == FrameworkError {
		return fmt.Sprintf("type : framework, code : %d, msg : %s",e.Code, e.Message)
	}
	return fmt.Sprintf("type : business, code : %d, msg : %s",e.Code, e.Message)
}

// new a framework type error
func NewFrameworkError(code int, msg string) *Error{
	return &Error{
		Type : FrameworkError,
		Code : code,
		Message : msg,
	}
}

// new a business type error
func New(code int, msg string) *Error{
	return &Error{
		Type : BusinuessError,
		Code : code,
		Message : msg,
	}
}
