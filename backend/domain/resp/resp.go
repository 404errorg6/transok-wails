package resp

type Err struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// NewErr creates a new error
func NewErr(code int, success bool, msg string) *Err {
	return &Err{
		Code:    code,
		Success: success,
		Msg:     msg,
	}
}

// Error implements the error interface
func (e *Err) Error() string {
	return e.Msg
}

// Out panics with the error function
func (e *Err) Out() {
	panic(*e)
}

// WithMsg sets the error message
func (e *Err) WithMsg(msg string) *Err {
	e.Msg = msg
	return e
}

// WithData sets the error data
func (e *Err) WithData(data interface{}) *Err {
	e.Data = data
	return e
}

/**
 * Predefined errors below
 */

// Success returns a success response
func Success() *Err {
	return NewErr(20000, true, "success")
}

// ClientErr returns a client error response
func ClientErr() *Err {
	return NewErr(40000, false, "client error")
}

// NotLoginErr returns an unauthorized error response
func NotLoginErr() *Err {
	return NewErr(40001, false, "unauthorized")
}

// Forbidden returns a forbidden error response
func Forbidden() *Err {
	return NewErr(40003, false, "forbidden")
}

// NotFoundErr returns a resource not found error response
func NotFoundErr() *Err {
	return NewErr(40004, false, "resource not found")
}

// NotExistErr returns a data not found error response
func NotExistErr() *Err {
	return NewErr(40005, false, "data not found")
}

// ExistErr returns a data already exists error response
func ExistErr() *Err {
	return NewErr(40006, false, "data already exists")
}

// InvalidErr returns an invalid data error response
func InvalidErr() *Err {
	return NewErr(40007, false, "invalid data")
}

// DataFormatErr returns a data format error response
func DataFormatErr() *Err {
	return NewErr(40010, false, "data format error")
}

// ServerErr returns a server error response
func ServerErr() *Err {
	return NewErr(50000, false, "server error")
}

// ServerBusyErr returns a server busy error response
func ServerBusyErr() *Err {
	return NewErr(50001, false, "server busy")
}

// ServerTimeoutErr returns a server timeout error response
func ServerTimeoutErr() *Err {
	return NewErr(50002, false, "server timeout")
}

// DbErr returns a database error response
func DbErr() *Err {
	return NewErr(50003, false, "database error")
}

// CacheErr returns a cache error response
func CacheErr() *Err {
	return NewErr(50004, false, "cache error")
}

// UnknownErr returns an unknown error response
func UnknownErr() *Err {
	return NewErr(99999, false, "unknown error")
}
