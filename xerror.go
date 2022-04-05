/* 
Package xerror provides additional functionailty to the standard errors library.

The standard errors library functionality is enhanced with external (user-facing)
error messages as well as internal (developer-facing) error codes.  A hexadecimal
error code is also included so that anyone experiencing an error can record the 
code which can be used to uniquely identify the specific error without relying
on copying down the verbose (and possibly duplicate) error message.

Finally, xerror also provides convenience methods to log the error message
and to perform standard HTTP error processing if the error occurs in an HTTP setting

*/
package xerror

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// Error is the interface that xerror implements
type Error interface {
	error
	Log()
	Handler(w http.ResponseWriter)
}

// Err is the base structure that includes all the information
// necessary for recording both External (user-facing) and 
// Internal (developer-facing) errors as well as the Code
// used as shorthand to refer to the error, the Status code
// which is used in the HTTP reply (if performed) and the 
// Location which indicates the source file and line number
// where the error ocurred.
type Err struct {
	Code     int			// a unique number representing this error
	Status   int			// HTTP status code for reply to HTTP client
	Location string   // source file and line number where error ocurred
	Internal error		// developer-oriented error message
	External error		// user-oriented error message 
}

// Errorf creates an xerror.Error given an internal error,
// a unique code and the format and components to include
// as the external error.  It resembles fmt.Errorf.
func Errorf(i error, c int, format string, a ...interface{}) *Err {
	// gets the source file and line where error ocurred
	// the runtime.Caller(1) is necessary so that the location
	// of the caller and not this code is recorded as the location
	_, file, line, _ := runtime.Caller(1)
	return &Err{
		Code:     c,
		Location: fmt.Sprintf("%s:%d", file, line),
		Internal: i,
		External: fmt.Errorf("error: "+format, a...),
	}
}

// Error provides the standard functionality whereby
// an error is rendered as a string.  In this case,
// the external (user-facing) error and error code
// are rendered.
func (e *Err) Error() string {
	return fmt.Sprintf("%s (0x%x)", e.External, e.Code)
}

// Log uses the standard log library to capture the
// internal error and location where the error ocurred
func (e *Err) Log() {
	log.Printf("error: %s", e.Internal)
	log.Printf("  at %s", e.Location)
}

// Handler performs the standard HTTP error handling
// functionality.  It also logs the error for later
// debugging, monitoring, etc.
func (e *Err) Handler(w http.ResponseWriter) {
	// use a default HTTP code if none provided
	if e.Status == 0 {
		e.Status = http.StatusInternalServerError
	}
	// record the external *and* internal errors
	// the naked "e" uses the Error() function
	log.Printf("%s: %s", e, e.Internal)
	log.Printf("  at %s", e.Location)
	http.Error(w, e.Error(), e.Status)
}
