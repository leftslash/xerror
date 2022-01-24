package xerror

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type Error interface {
	error
	Log()
	Handler(w http.ResponseWriter)
}

type Err struct {
	Code     int
	Status   int
	Location string
	Internal error
	External error
}

func Errorf(i error, c int, format string, a ...interface{}) *Err {
	_, file, line, _ := runtime.Caller(1)
	return &Err{
		Code:     c,
		Location: fmt.Sprintf("%s:%d", file, line),
		Internal: i,
		External: fmt.Errorf("error: "+format, a...),
	}
}

func (e *Err) Error() string {
	return fmt.Sprintf("%s (%x)", e.External, e.Code)
}

func (e *Err) Log() {
	log.Printf("error: %s", e.Internal)
	log.Printf("  at %s", e.Location)
}

func (e *Err) Handler(w http.ResponseWriter) {
	if e.Status == 0 {
		e.Status = http.StatusInternalServerError
	}
	log.Printf("%s: %s", e, e.Internal)
	log.Printf("  at %s", e.Location)
	http.Error(w, e.Error(), e.Status)
}
